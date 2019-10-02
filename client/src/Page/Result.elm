port module Page.Result exposing (Model, Msg, init, subscriptions, update, view)

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
  { title = "test | result"
  , body =
    [ div [] [] ]
  }