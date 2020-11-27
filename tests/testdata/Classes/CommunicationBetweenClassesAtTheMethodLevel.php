<?php

class A {
  public function AMethod() {
    $b = new B; // ok
    $b->BMethod(); // ok
    C::CMethod(); // ok
    $d = new D; // ok
  }
}

class B {
  public function BMethod() {}
}

class C {
  public static function CMethod() {}
}

class D {
  public function DMethod() {}
}

class E {
  public function EMethod() {}
}
