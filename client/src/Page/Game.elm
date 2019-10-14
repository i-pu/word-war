port module Page.Game exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

port startGame : String -> Cmd msg
port say : Message -> Cmd msg
port onMessage : ( Message -> msg ) -> Sub msg

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
  ( Model env "" [] { userId = "" }
  , startGame "testuid"
  )

type Msg
  = MessageInputChange String
  | StartGame
  | Say
  | OnMessage Message

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

subscriptions : Model -> Sub Msg
subscriptions _ =
  onMessage OnMessage

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | game"
  , body =
    [ div [ class "container" ]
      [ h3 [] [ text "Messages" ]
      , viewMessages model.messages
      , input [ type_ "text"
        , onInput MessageInputChange
        , value model.messageInput
        ] []
      , button [ onClick Say ] [ text "Send" ]
      , a [ Route.href <| Route.Result ] [ text "/result" ]
      ]
    ]
  }

hero : Html Msg
hero =
  section [ class "hero is-primary" ]
    [ div [ class "hero-body" ]
      [ div [ class "container" ]
        [ h1 [ class "title" ]
          [ text "Game" ]
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
