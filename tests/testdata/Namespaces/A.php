<?php

namespace A {
  class AClass {
    public static function AMethod() {
      \B\BClass::BMethod();
      \D\DClass::DMethod();
    }
  }

  function AFunc() {
    \C\CClass::CMethod();
    \E\EFunc();
  }
}

namespace B {
  class BClass {
    public static function BMethod() {
      \C\CClass::CMethod();
    }
  }
}

namespace C {
  class CClass {
    public static function CMethod() {
      \A\AClass::AMethod();
    }
  }
}

namespace D {
  class DClass {
    public static function DMethod() {
      \A\AFunc();
    }
  }
}

namespace E {
  function EFunc() {
    \C\CClass::CMethod();
    GlobalFunction();
  }
}
