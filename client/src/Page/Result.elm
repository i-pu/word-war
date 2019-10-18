port module Page.Result exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env, navKey, getUid)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Route
import String exposing (fromInt)

-- to js
port requestResult : ( String ) -> Cmd msg

-- from js
port onResult : ( Result -> msg ) -> Sub msg

type alias Result =
  { userId : String
  , score : Int
  }

type alias Model =
  { env : Env
  , result : Result
  }

init : Env -> ( Model, Cmd Msg )
init env =
  ( Model env { userId = "", score = 0 }
  , requestResult (getUid env)
  )

type Msg
  = OnResult Result
  | ToHome

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    OnResult newResult ->
      ({ model | result = newResult }, Cmd.none)
    ToHome ->
      (model, Route.replaceUrl (navKey model.env) Route.Home)

subscriptions : Model -> Sub Msg
subscriptions _ =
  Sub.batch 
  [ onResult OnResult
  ]

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | result"
  , body =
    [ hero
    , section [ class "section" ] [ resultView model ]
    ]
  }

hero : Html Msg
hero =
  section [ class "hero is-primary" ]
    [ div [ class "hero-body" ]
      [ div [ class "container" ]
        [ h1 [ class "title" ]
          [ text "Result" ]
        ]
      ]
    ]

resultView : Model -> Html Msg
resultView model =
  div [ class "container has-text-centered" ]
    [ h1 [ class "title" ] [ text (fromInt model.result.score ++ "点") ]
    , button [ class "button is-primary is-large is-fullwidth", onClick ToHome ] [ text "ホームへ" ]
    ]