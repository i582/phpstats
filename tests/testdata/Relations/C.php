<?php

class SomeClass {
  public function Method() {
    $this->OtherMethod();
  }

  public function OtherMethod() {
    someFunc();
  }
}

function someFunc() {

}

function someOtherFunc() {
  $s = new SomeClass();
  $s->Method();
}
