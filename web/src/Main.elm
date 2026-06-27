module Main exposing (..)

import Element exposing (..)
import Element.Background as Background
import Element.Font as Font
import Element.Border as Border
import Element.Input as Input

import Browser
import Browser.Events
import Browser.Navigation
import Browser.Dom
import Html.Attributes

import Task
import Url


main : Program () Model Msg
main =
  Browser.application
  { init = init
  , update = update
  , view = boneView
  , subscriptions = subscriptions
  , onUrlChange = onUrlChange
  , onUrlRequest = onUrlRequest
  }

type alias Model =
  { height : Int
  }

type Msg = None
  | Resize Int Int
  | InitialResize Browser.Dom.Viewport


onUrlChange : Url.Url -> Msg
onUrlChange _ = None

onUrlRequest : Browser.UrlRequest -> Msg
onUrlRequest _ = None


init : () -> Url.Url -> Browser.Navigation.Key -> ( Model, Cmd Msg )
init _ _ _ =
  ( { height = 200
    }
  , Task.perform InitialResize Browser.Dom.getViewport
  )

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model = case msg of
  None ->
    ( model, Cmd.none )
  Resize _ h ->
    ( {model | height = h}, Cmd.none)
  InitialResize viewport ->
    ( { model | height = truncate viewport.viewport.height }, Cmd.none )

boneView : Model -> Browser.Document Msg
boneView m =
  { title = "hey"
  , body = [ layout globalStyles (view m) ]
  }

view : Model -> Element Msg
view m = row
  [ width fill
  , height (px m.height)
  ]
  [ leftBar m
  , mailInbox m
  ]

leftBar : Model -> Element Msg
leftBar m = column
  [ width (fillPortion 1)
  , padding 40
  , alignTop
  ]
  [ (title "Seamail")
  , controls m
  ]

type alias Mail =
  { id : Int
  , address : String
  , subject : String
  }

mailTest : List Mail
mailTest =
  List.repeat 40
    { id = 1
    , address = "okay@gmail.com"
    , subject = "How a re you ?"}

mailInbox : Model -> Element Msg
mailInbox _ = column
  [ width (fillPortion 3)
  , alignTop
  , scrollbarY
  , height fill
  ]
  ( List.map mailShowcase mailTest )

mailShowcase : Mail -> Element Msg
mailShowcase mail = link
  [ width fill ]
  { url="/", label=row
    [ width fill
    , Background.color palette.surface
    , Border.color palette.surface2
    , Border.width 1
    , mouseOver
      [ Background.color palette.surface2
      ]
    , padding 5
    , spacing 24
    , cursorPointer
    ]
    [ Input.checkbox
      [ width (fillPortion 1)
      , padding 5
      ]
      { onChange = \_ -> None
      , checked = False
      , icon = Input.defaultCheckbox
      , label = Input.labelHidden "Checkbox"
      }
    , el
      [ width (fillPortion 4)
      ] ( text mail.address )
    , el
      [ width (fillPortion 12)
      , Font.color palette.mutedText
      ]
      ( text mail.subject )
    ]
  }


type alias Palette =
    { bg : Color
    , surface : Color
    , surface2 : Color
    , border : Color
    , mutedText : Color
    , text : Color
    }

palette : Palette
palette =
    { bg = rgb255 10 10 10
    , surface = rgb255 23 23 23
    , surface2 = rgb255 38 38 38
    , border = rgb255 64 64 64
    , mutedText = rgb255 163 163 163
    , text = rgb255 250 250 250
    }


globalStyles : List (Attr () msg)
globalStyles =
  [ Background.color palette.bg
  , Font.color palette.text
  , width fill
  ]


title : String -> Element msg
title s = link
  [ Font.size 40
  , Font.bold
  , cursorPointer
  , paddingXY 0 20
  ]
  { url = "/"
  , label = text s
  }

controls : Model -> Element msg
controls _ = column
  [ width fill
  ]
  [ action "Inbox"
  , action "Send"
  ]

action : String -> Element msg
action s = link
  [ mouseOver [ Background.color palette.surface2 ]
  , Border.rounded 5
  , width (maximum 350 fill)
  , padding 10
  ]
  { url = "/"
  , label = row
    []
    [ (text s)
    ]
  }


cursorPointer : Attribute msg
cursorPointer = htmlAttribute (Html.Attributes.style "cursor" "pointer")

subscriptions _ =
  Browser.Events.onResize Resize
