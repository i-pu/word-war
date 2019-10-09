port module Page.Result exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

-- to js
port requestResult : ( String ) -> Cmd msg

-- from js
port resultCallback : ( Result -> msg ) -> Sub msg

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
  , Cmd.none
  )

type Msg
  = OnResult Result
  | Request

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    OnResult newResult ->
      ({ model | result = newResult }, Cmd.none)
    Request ->
      (model, requestResult "uid")

subscriptions : Model -> Sub Msg
subscriptions _ =
  Sub.batch 
  [ resultCallback OnResult
  ]

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | result"
  , body =
    [ hero
    , a [ Route.href <| Route.Home ] [ text "/home" ]
    , button [ onClick Request ] [ text "push me" ]
    , h3 [] [ text model.result.userId ]
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