port module Page.Home exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

type alias Message =
  { name : String
  , message : String
  }

type alias Model =
  { env : Env
  }

init : Env -> ( Model, Cmd Msg )
init env =
  ( Model env
  , Cmd.none
  )

type Msg
  = Hoge

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    Hoge ->
      (model, Cmd.none)

subscriptions : Model -> Sub Msg
subscriptions _ =
  Sub.none

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | home"
  , body =
    [ hero
    , a [ Route.href <| Route.Game ] [ text "/game" ]
    ]
  }

hero : Html Msg
hero =
  section [ class "hero is-primary" ]
    [ div [ class "hero-body" ]
      [ div [ class "container" ]
        [ h1 [ class "title" ]
          [ text "Home" ]
        ]
      ]
    ]