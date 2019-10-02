port module Page.Game exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

port toJS : Message -> Cmd msg
port toElm : ( Message -> msg ) -> Sub msg

type alias Message =
  { name : String
  , message : String
  }

type alias Model =
  { env : Env
  , messageInput : String
  , messages : List Message
  }

init : Env -> ( Model, Cmd Msg )
init env =
  ( Model env "" []
  , Cmd.none
  )

type Msg
  = MessageInputChange String
  | SendToJS
  | NewMessage Message

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    MessageInputChange newInput ->
      ({ model | messageInput = newInput }, Cmd.none)
    SendToJS ->
      ({ model | messageInput = "" }, toJS ({ name = "hoge", message = model.messageInput}))
    NewMessage incoming ->
      ({ model | messages = model.messages ++ [ incoming ] }, Cmd.none)

subscriptions : Model -> Sub Msg
subscriptions _ =
  toElm NewMessage

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
      , button [ onClick SendToJS ] [ text "Send" ]
      ]
    ]
  }

viewMessages : List Message -> Html Msg
viewMessages messages =
  ul []
    (List.map viewMessage messages)

viewMessage : Message -> Html Msg
viewMessage message =
  li []
    [ text (message.name ++ ":" ++ message.message) ]
    