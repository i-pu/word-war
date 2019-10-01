port module Page.Index exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

-- to JS
port toJS : Message -> Cmd msg
port signupWithFirebase : ({ email : String, password : String }) -> Cmd msg
port signinWithFirebase : ({ email : String, password : String }) -> Cmd msg

-- from JS
port toElm : ( Message -> msg ) -> Sub msg
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
  , messageInput : String
  , messages : List Message
  }

init : Env -> ( Model, Cmd Msg )
init env =
  ( Model env "" "" "" "" []
  , Cmd.none
  )

type Msg
  = MessageInputChange String
  | EmailInputChange String
  | PasswordInputChange String
  | OnClickSignUpButton
  | OnClickSignInButton
  | SendToJS
  | NewMessage Message
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
    MessageInputChange newInput ->
      ({ model | messageInput = newInput }, Cmd.none)
    SendToJS ->
      ({ model | messageInput = "" }, toJS ({ name = "hoge", message = model.messageInput}))
    NewMessage incoming ->
      ({ model | messages = model.messages ++ [ incoming ] }, Cmd.none)
    OnSignin user ->
      ({ model | uid = user.uid }, Cmd.none)

subscriptions : Model -> Sub Msg
subscriptions _ =
  Sub.batch
    [ toElm NewMessage
    , signinCallback OnSignin
    ]

view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | infex"
  , body =
    [ h3 [] [ text "Signup with firebase/auth" ]
    , signupForm model
    , h3 [] [ text "Messages" ]
    , viewMessages model.messages
    , input [ type_ "text"
      , onInput MessageInputChange
      , value model.messageInput
      ] []
    , button [ onClick SendToJS ] [ text "Send" ]
    ]
  }

signupForm : Model -> Html Msg
signupForm model =
  div [] 
    [ h4 [] [ text ("uid: " ++ model.uid) ]
    , input [ type_ "text"
      , placeholder "email"
      , onInput EmailInputChange
      , value model.emailInput
      ] []
    , br [] []
    , input [ type_ "password"
      , placeholder "password"
      , onInput PasswordInputChange
      , value model.passwordInput
      ] []
    , br [] []
    , button [ onClick OnClickSignInButton ] [ text "Sign in" ]
    , button [ onClick OnClickSignUpButton ] [ text "Sign up" ]
  ]

viewMessages : List Message -> Html Msg
viewMessages messages =
  ul []
    (List.map viewMessage messages)

viewMessage : Message -> Html Msg
viewMessage message =
  li []
    [ text (message.name ++ ":" ++ message.message) ]