port module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)

port toJS : String -> Cmd msg
port toElm : ( String -> msg ) -> Sub msg

main : Program String Model Msg
main =
  Browser.element
    { init = init
    , update = update
    , view = view
    , subscriptions = subscriptions
    }

type alias Model =
  { input : String
  , messages : List String
  }

type Msg
  = InputChange String
  | SendToJS
  | NewMessage String

init : String -> ( Model, Cmd Msg )
init flags =
  (
    { input = ""
    , messages = []
    }, 
    Cmd.none
  )

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    InputChange newInput ->
      ({ model | input = newInput }, Cmd.none)
    SendToJS ->
      ({ model | input = "" }, toJS model.input)
    NewMessage incoming ->
      ({ model | messages = model.messages ++ [ incoming ] }, Cmd.none)

subscriptions : Model -> Sub Msg
subscriptions model =
  toElm NewMessage

view : Model -> Html Msg
view model =
  div []
    [ viewMessages model.messages
    , input [ type_ "text"
      , onInput InputChange
      , value model.input
      ] []
    , button [ onClick SendToJS ] [ text "Send" ]
    ]

viewMessages : List String -> Html Msg
viewMessages messages =
  ul []
    (List.map viewMessage messages)

viewMessage : String -> Html Msg
viewMessage message =
  li []
    [ text message ]