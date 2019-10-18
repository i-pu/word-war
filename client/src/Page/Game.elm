port module Page.Game exposing (Model, Msg, init, subscriptions, update, view, isHiragana)

import Env exposing (Env, navKey)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Regex
import Id exposing (Id)
import Route

port startGame : String -> Cmd msg
port say : Message -> Cmd msg

port onMessage : ( Message -> msg ) -> Sub msg
port onFinish : ( () -> msg ) -> Sub msg

type alias Message =
  { userId : String
  , message : String
  }

type alias Model =
  { env : Env
  , messageInput : String
  , messages : List Message
  , user : { userId : String }
  }

init : Env -> ( Model, Cmd Msg )
init env =
  ( Model env "" [] { userId = "testuid" }
  , startGame "testuid"
  )

type Msg
  = MessageInputChange String
  | StartGame
  | Say
  | OnMessage Message
  | OnFinish ()

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  let
    user = model.user
  in
    case msg of
      MessageInputChange newInput ->
        ({ model | messageInput = newInput }, Cmd.none)
      StartGame ->
        ( model, startGame user.userId )
      Say ->
        ({ model | messageInput = "" }, say ({ userId = user.userId, message = model.messageInput}))
      OnMessage incoming ->
        ({ model | messages = model.messages ++ [ incoming ] }, Cmd.none)
      OnFinish _ ->
        ( model, Route.replaceUrl (navKey model.env) Route.Result )

subscriptions : Model -> Sub Msg
subscriptions _ = Sub.batch
  [ onMessage OnMessage
  , onFinish OnFinish
  ]

hiragana : Regex.Regex
hiragana =
  Maybe.withDefault Regex.never <| Regex.fromString "^[ぁ-んー]+$"

isHiragana : String -> Bool
isHiragana str =
  Regex.find hiragana str
    |> List.isEmpty
    |> not

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | game"
  , body =
    [ section [ class "section" ]
      [ form model
      , viewMessages model.messages
      ]
    ]
  }

form : Model -> Html Msg
form model =
  div [ class "container" ]
    [ div [ class "field has-addons" ]
      [ div [ class "control is-expanded" ]
        [ input
          [ type_ "text"
          , class "input"
          , placeholder "Word"
          , onInput MessageInputChange
          , value model.messageInput
          ] []
        ]
      , div [ class "control" ]
        [ button 
          [ class "button is-primary"
          , onClick Say
          , disabled <| not (isHiragana model.messageInput)
          ]
          [ text "Send" ]
        ]
      ]
    ]

viewMessages : List Message -> Html Msg
viewMessages messages =
  ul []
    (List.map viewMessage messages)

viewMessage : Message -> Html Msg
viewMessage message =
  li []
    [ text (message.userId ++ ":" ++ message.message) ]
