module Route exposing (Route(..), fromUrl, href, parser, routeToString)

import Html exposing (Attribute)
import Html.Attributes as Attr
import Id exposing (Id)
import Url exposing (Url)
import Env exposing (Env)
import Url.Parser as Parser exposing ((</>), Parser)

type Route
    = Top
    | Home
    | Game
    | Result

parser : Parser (Route -> a) a
parser =
    Parser.oneOf
        [ Parser.map Top Parser.top
        , Parser.map Home (Parser.s "view")
        , Parser.map Game (Parser.s "game")
        , Parser.map Result (Parser.s "result")
        ]

fromUrl : Url -> Maybe Route
fromUrl url =
    Parser.parse parser url

href : Route -> Attribute msg
href targetRoute =
    Attr.href (routeToString targetRoute)

routeToString : Route -> String
routeToString page =
    let
        pieces =
            case page of
                Top ->
                    []
                Home ->
                    [ "view" ]
                Game ->
                    [ "game" ]
                Result ->
                    [ "result" ]
    in
    String.join "/" pieces