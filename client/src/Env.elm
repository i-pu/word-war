module Env exposing (Env, User, create, navKey)

import Browser.Navigation as Nav

-- ユーザーモデル
type alias User =
  { uid : String
  }

type Env
    = Env Internals

-- [TODO] セッションデータをもたせたい
type alias Internals =
    { key : Nav.Key
    -- , user : User
    }

create : Nav.Key -> Env
create key =
    Env (Internals key)

navKey : Env -> Nav.Key
navKey (Env internals) =
    internals.key
