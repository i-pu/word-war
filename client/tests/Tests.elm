module Tests exposing (..)
import Test.Html.Query as Query
import Test exposing (..)
import Test.Html.Selector exposing (tag, text)
import Html exposing (..)

import Page.Top exposing (hero)

-- Check out https://package.elm-lang.org/packages/elm-explorations/test/latest to learn more about testing in Elm!

-- all : Test
-- all =
--     describe "A Test Suite"
--         [ test "Addition" <|
--             \_ ->
--                 Expect.equal 10 (3 + 7)
--         , test "String.left" <|
--             \_ ->
--                 Expect.equal "a" (String.left 1 "abcdefg")
--         , test "This test should fail" <|
--             \_ ->
--                 Expect.fail "failed as expected!"
--         ]

all : Test
all = 
    describe "Top snapshot test"
        [ test "Form" <|
            \_ ->
                hero
                |> Query.fromHtml
                |> Query.has [ tag "h1" ]
        ]
