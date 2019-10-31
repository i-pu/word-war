module Env exposing (Env, create, navKey, setUid, getUid)

import Browser.Navigation as Nav

type Env
  = Env Internals

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

getUid : Env -> String
getUid (Env internals) =
  internals.uid

setUid : Env -> String -> Env
setUid (Env internals) uid =
  Env { internals | uid = uid}
