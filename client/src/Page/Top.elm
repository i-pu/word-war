port module Page.Top exposing (Model, Msg, init, subscriptions, update, view)

import Env exposing (Env, User, navKey)
import Html exposing (..)
import Html.Events exposing (..)
import Html.Attributes exposing (..)
import Id exposing (Id)
import Route

-- 機能の追加の仕方
-- 1. type msg に update で使いたいイベントを増やす
-- 2. update の case にモデルと関数を呼ぶ処理を増やす
-- 3. view と連携したいのであれば 要素の中の onClickなどにtype msg にあるイベントを書く

-- JS を呼ぶ関数をここに書く
port signupWithFirebase : ({ email : String, password : String }) -> Cmd msg
port signinWithFirebase : ({ email : String, password : String }) -> Cmd msg

-- JS から呼ばれる関数をここに書いておく
port signinCallback : ( User -> msg) -> Sub msg

-- メッセージモデル
type alias Message =
  { name : String
  , message : String
  }

-- モデル(画面で扱うデータ)
type alias Model =
  { env : Env -- 全体で共有する環境
  , emailInput : String -- Emailのテキストボックスの中身
  , passwordInput : String -- Passwordのテキストボックスの中身
  }

init : Env -> ( Model, Cmd Msg )
init env =
  -- モデルを初期化
  ( Model env "" ""
  , Cmd.none
  )

-- アクションのイベントを列挙する
type Msg
  = EmailInputChange String
  | PasswordInputChange String
  | OnClickSignUpButton
  | OnClickSignInButton
  | OnSignin User

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
  case msg of
    -- Email テキストボックスの中身が変更された時
    EmailInputChange newInput ->
      ({ model | emailInput = newInput }, Cmd.none)
    -- Password テキストボックスの中身が変更された時
    PasswordInputChange newInput ->
      ({ model | passwordInput = newInput }, Cmd.none)
    -- SignUp ボタンがクリックされた時
    OnClickSignUpButton ->
      (
        -- テキストボックスの中身を空に
        { model 
        | emailInput = ""
        , passwordInput = "" 
        },
        -- JS の 関数を呼ぶ
        signupWithFirebase (
          { email = model.emailInput
          , password = model.passwordInput
          }
        )
      )
    -- SignIn ボタンがクリックされた時
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
      (model
      , Cmd.batch 
        [ Route.replaceUrl (navKey model.env) Route.Home
        ]
      )

-- JS から呼ばれる関数を登録する
subscriptions : Model -> Sub Msg
subscriptions _ =
  -- 関数が複数になるときは Sub.batch を使う
  Sub.batch
    [ signinCallback OnSignin
    ]

-- HTML にあたる部分
-- 要素 [ アトリビュート ] [ 中身 ] の形で

-- 例 <h1 class="hoge">Hello</h1>
-- は Elm では h1 [ class "hoge" ] [ text "Hello" ] になる
view : Model -> { title : String, body : List (Html Msg) }
view model =
  { title = "test | top"
  , body =
    -- hero は 関数として切り分けてコンポーネント化している
    [ hero
    , div [ class "container" ] (signupForm model)
    ]
  }

-- 返り値の型を Html Msg にすると Html の要素を返せる
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

-- 入力フォームを返す関数
signupForm : Model -> List (Html Msg)
signupForm model =
  [ div [ class "field" ]
    [ label [ class "label" ] [ text "Username" ]
    , div [ class "control has-icons-left has-icons-right" ]
      [ input [ class "input is-success", type_ "text", placeholder "Email", value model.emailInput, onInput EmailInputChange ]
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