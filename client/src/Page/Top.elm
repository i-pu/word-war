port module Page.Top exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

-- to JS
port signupWithFirebase : ({ email : String, password : String }) -> Cmd msg
port signinWithFirebase : ({ email : String, password : String }) -> Cmd msg

-- from JS
port signinCallback : ( User -> msg) -> Sub msg

type alias User =
  { uid : String
  }

type alias Message =
  { name : String
  , message : String
  }

type alias Model =
  { env : Env
  , uid : String
  , emailInput : String
  , passwordInput : String
  }

init : Env -> ( Model, Cmd Msg )
init env =
  ( Model env "" "" ""
  , Cmd.none
  )

type Msg
  = EmailInputChange String
  | PasswordInputChange String
  | OnClickSignUpButton
  | OnClickSignInButton
  | OnSignin User

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    EmailInputChange newInput ->
      ({ model | emailInput = newInput }, Cmd.none)
    PasswordInputChange newInput ->
      ({ model | passwordInput = newInput }, Cmd.none)
    OnClickSignUpButton ->
      (
        { model 
        | emailInput = ""
        , passwordInput = "" 
        },
        signupWithFirebase (
          { email = model.emailInput
          , password = model.passwordInput
          }
        )
      )
    OnClickSignInButton ->
      (
        { model 
        | emailInput = ""
        , passwordInput = "" 
        },
        signinWithFirebase (
          { email = model.emailInput
          , password = model.passwordInput
          }
        )
      )
    OnSignin user ->
      ({ model | uid = user.uid }, Cmd.none)

subscriptions : Model -> Sub Msg
subscriptions _ =
  Sub.batch
    [ signinCallback OnSignin
    ]

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | top"
  , body =
    [ hero
    , div [ class "container" ] (signupForm model)
    ]
  }

hero : Html Msg
hero =
  section [ class "hero is-primary" ]
    [ div [ class "hero-body" ]
      [ div [ class "container" ]
        [ h1 [ class "title" ]
          [ text "Word War" ]
        , h2 [ class "subtitle" ]
          [ text "リアルタイムしりとり" ]
        ]
      ]
    ]

signupForm : Model -> List (Html Msg)
signupForm model =
  [ div [ class "field" ]
    [ label [ class "label" ] [ text "Username" ]
    , div [ class "control has-icons-left has-icons-right" ]
      [ input [ class "input is-success", type_ "text", placeholder "Text input", value model.emailInput, onInput EmailInputChange ]
        []
      , span [ class "icon is-small is-left" ]
        [ i [ class "fas fa-envelope" ] 
          []
        ]
      ]
    ]
  , div [ class "field" ]
    [ label [ class "label" ] [ text "Password" ]
    , div [ class "control has-icons-left" ]
      [ input [ class "input", type_ "password", placeholder "Password", value model.passwordInput, onInput PasswordInputChange ]
        []
      , span [ class "icon is-small is-left" ]
        [ i [ class "fas fa-lock" ]
          []
        ]
      ]
    ]
  , div [ class "field is-grouped" ]
    [ div [ class "control" ]
      [ button [ class "button is-success", onClick OnClickSignInButton ]
        [ text "Signin" ]
      ]
    , div [ class "control" ]
      [ button [ class "button", onClick OnClickSignUpButton ]
        [ text "Signup" ]
      ]
    ]
  ]