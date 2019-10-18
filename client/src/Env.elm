module Env exposing (Env, User, create, navKey, setUid, getUid)

import Browser.Navigation as Nav

type Env
  = Env Internals

type alias User =
  { uid : String
  }

type alias Internals =
  { key : Nav.Key
  , uid : String
  }

create : Nav.Key -> String -> Env
create key str =
  Env (Internals key str)

navKey : Env -> Nav.Key
navKey (Env internals) =
  internals.key

-- [TODO] 保存方法考える

getUid : Env -> String
getUid (Env internals) =
  "ESHNsCBaZHVQEjgDeEB2VXSCG262"
  -- internals.uid

setUid : Env -> String -> Env
setUid (Env internals) uid =
  create internals.key uid
