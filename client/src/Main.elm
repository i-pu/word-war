module Main exposing (Flags, Model(..), Msg(..), changeRouteTo, init, main, subscriptions, toEnv, update, updateWith, view)

import Browser
import Browser.Navigation as Nav
import Env exposing (Env)
import Html exposing (Html, button, div, text)
import Html.Events exposing (onClick)
import Id exposing (Id)
import Page.Top as TopPage
import Page.Home as HomePage
import Page.Game as GamePage
import Page.Result as ResultPage
import Route exposing (Route)
import Url

main : Program Flags Model Msg
main =
    Browser.application
        { init = init
        , update = update
        , subscriptions = subscriptions
        , view = view
        , onUrlChange = UrlChanged
        , onUrlRequest = LinkClicked
        }

-- MODEL
type Model
    = NotFound Env
    | Top Env TopPage.Model
    | Home Env HomePage.Model
    | Game Env GamePage.Model
    | Result Env ResultPage.Model

type alias Flags =
    {}

init : Flags -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
    changeRouteTo (Route.fromUrl url)
        (NotFound <|
            Env.create key
        )

-- UPDATE
type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | GotTopMsg TopPage.Msg
    | GotHomeMsg HomePage.Msg
    | GotGameMsg GamePage.Msg
    | GotResultMsg ResultPage.Msg

update : Msg -> Model -> ( Model, Cmd Msg )
update message model =
    let
        env =
            toEnv model
    in
    case ( message, model ) of
        ( LinkClicked urlRequest, _ ) ->
            case urlRequest of
                Browser.Internal url ->
                    case Route.fromUrl url of
                        Just _ ->
                            ( model, Nav.pushUrl (Env.navKey env) (Url.toString url) )

                        Nothing ->
                            ( model, Nav.load <| Url.toString url )

                Browser.External href ->
                    if String.length href == 0 then
                        ( model, Cmd.none )

                    else
                        ( model, Nav.load href )

        ( UrlChanged url, _ ) ->
            changeRouteTo (Route.fromUrl url) model

        ( GotTopMsg subMsg, Top _ subModel ) ->
            TopPage.update subMsg subModel
                |> updateWith (Top env) GotTopMsg

        ( GotHomeMsg subMsg, Home _ subModel ) ->
            HomePage.update subMsg subModel
                |> updateWith (Home env) GotHomeMsg

        ( GotGameMsg subMsg, Game _ subModel ) ->
            GamePage.update subMsg subModel
                |> updateWith (Game env) GotGameMsg
        
        ( GotResultMsg subMsg, Result _ subModel ) ->
            ResultPage.update subMsg subModel
                |> updateWith (Result env) GotResultMsg

        ( _, _ ) ->
            ( model, Cmd.none )


changeRouteTo : Maybe Route -> Model -> ( Model, Cmd Msg )
changeRouteTo maybeRoute model =
    let
        env =
            toEnv model
    in
    case maybeRoute of
        Just Route.Top ->
            TopPage.init env
                |> updateWith (Top env) GotTopMsg

        Just Route.Home ->
            HomePage.init env
                |> updateWith (Home env) GotHomeMsg

        Just Route.Game ->
            GamePage.init env
                |> updateWith (Game env) GotGameMsg

        Just Route.Result ->
            ResultPage.init env
                |> updateWith (Result env) GotResultMsg

        Nothing ->
            ( NotFound env, Cmd.none )

        -- Just (Route.View id) ->
        --     ViewPage.init env id
        --         |> updateWith (View env id) GotViewMsg

toEnv : Model -> Env
toEnv page =
    case page of
        NotFound env ->
            env

        Top env _ ->
            env
        
        Home env _ ->
            env

        Game env _ ->
            env

        Result env _ ->
            env

        -- View env _ _ ->
        --     env


updateWith : (subModel -> Model) -> (subMsg -> Msg) -> ( subModel, Cmd subMsg ) -> ( Model, Cmd Msg )
updateWith toModel toMsg ( subModel, subCmd ) =
    ( toModel subModel
    , Cmd.map toMsg subCmd
    )

-- SUBSCRIPTIONS
subscriptions : Model -> Sub Msg
subscriptions model =
    case model of
        NotFound _ ->
            Sub.none

        Top _ subModel ->
            Sub.map GotTopMsg (TopPage.subscriptions subModel)

        Home _ subModel ->
            Sub.map GotHomeMsg (HomePage.subscriptions subModel)

        Game _ subModel ->
            Sub.map GotGameMsg (GamePage.subscriptions subModel)

        Result _ subModel ->
            Sub.map GotResultMsg (ResultPage.subscriptions subModel)

        -- View _ _ subModel ->
        --     Sub.map GotViewMsg (ViewPage.subscriptions subModel)

-- VIEW


view : Model -> Browser.Document Msg
view model =
    let
        viewPage toMsg { title, body } =
            { title = title, body = List.map (Html.map toMsg) body }
    in
    case model of
        NotFound _ ->
            { title = "Not Found", body = [ Html.text "Not Found" ] }

        Top _ subModel ->
            viewPage GotTopMsg (TopPage.view subModel)

        Home _ subModel ->
            viewPage GotHomeMsg (HomePage.view subModel)

        Game _ subModel ->
            viewPage GotGameMsg (GamePage.view subModel)

        Result _ subModel ->
            viewPage GotResultMsg (ResultPage.view subModel)

        -- View _ _ subModel ->
        --     viewPage GotViewMsg (ViewPage.view subModel)
        