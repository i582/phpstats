<?php

class LCOM {
  const CONSTANT_ZERO   = 0;
  const CONSTANT_ONE    = 1;
  const CONSTANT_UNUSED = -1;

  public int $publicProp;
  private int $privateProp;
  public int $unusedProp;
  public int $group1Prop;

  public function publicMethod() {
    LCOM::CONSTANT_ZERO;
    $this->internalMethod();
    $this->publicProp = 2;
    $this->privateProp = 10;
  }

  private function internalMethod() {
    LCOM::CONSTANT_ONE;
    $this->privateProp = 3;
  }

  public function unusedMethod() {
  }

  public function group1Method() {
    $this->group1Prop = 10;
  }
}
