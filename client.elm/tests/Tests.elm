module Tests exposing (..)
import Test.Html.Query as Query
import Test exposing (..)
import Test.Html.Selector exposing (tag, text)
import Html exposing (..)
import Expect

import Page.Top exposing (init, hero)
import Page.Game exposing (isHiragana)

-- Check out https://package.elm-lang.org/packages/elm-explorations/test/latest to learn more about testing in Elm!

-- all : Test
-- all =
--   describe "A Test Suite"
--     [ test "Addition" <|
--       \_ ->
--         Expect.equal 10 (3 + 7)
--     , test "String.left" <|
--       \_ ->
--         Expect.equal "a" (String.left 1 "abcdefg")
--     , test "This test should fail" <|
--       \_ ->
--         Expect.fail "failed as expected!"
--     ]

-- want to include Browser.Navigation in testing
-- <https://github.com/elm-explorations/test/issues/24>
topTest : Test
topTest = 
  describe "Top test"
    [ test "Hero" <|
      \_ ->
        hero
        |> Query.fromHtml
        |> Query.has [ tag "h1" ]
    ]

gameTest : Test
gameTest =
  describe "Game test"
    [ test "isHiragana1" <|
      \_ ->
        Expect.equal True (isHiragana "あい")
    , test "isHiragana2" <|
      \_ ->
        Expect.equal False (isHiragana "aff")
    , test "isHiragana3" <|
      \_ ->
        Expect.equal False (isHiragana "あaff")
    ]

