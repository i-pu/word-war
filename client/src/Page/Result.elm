port module Page.Result exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
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
  , requestResult "testuid"
  )

type Msg
  = OnResult Result

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    OnResult newResult ->
      ({ model | result = newResult }, Cmd.none)

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
    , h3 [] [ text (fromInt model.result.score ++ "点") ]
    , a [ Route.href <| Route.Home ] [ text "ホームへ" ]
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