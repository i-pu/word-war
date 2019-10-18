port module Page.Home exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env, navKey)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

type alias Model =
  { env : Env
  }

init : Env -> ( Model, Cmd Msg )
init env =
  ( Model env
  , Cmd.none
  )

type Msg
  = ToGame

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    ToGame ->
      (model, Route.replaceUrl (navKey model.env) Route.Game)

subscriptions : Model -> Sub Msg
subscriptions _ =
  Sub.none

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | home"
  , body =
    [ hero
    , selectMenu
    ]
  }

hero : Html Msg
hero =
  section [ class "hero is-primary" ]
    [ div [ class "hero-body" ]
      [ div [ class "container" ]
        [ h1 [ class "title" ]
          [ text ("{Name}" ++ "Rating: 0.00") ]
        ]
      ]
    ]

selectMenu : Html Msg
selectMenu =
  section [ class "section" ]
    [ div [ class "container" ]
      [ button [ style "margin-bottom" "10px", class "button is-primary is-large is-fullwidth", onClick ToGame ] [ text "ランダムマッチ" ]
      , button [ class "button is-link is-large is-fullwidth", disabled True ] [ text "部屋を作る" ]
      ]
    ]