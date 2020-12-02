<?php

class TargetClassA {
  const CONSTANT = 0;

  public int $field = 10;

  public static function targetMethod1() {
    TargetClassB::targetMethod();
  }

  public static function targetMethod2() {
    TargetClassB::targetMethod();
  }
}

class TargetClassB {
  public static function targetMethod() {
    TargetClassA::targetMethod();
    $tga = new TargetClassA();
    $tga->field;
    TargetClassA::CONSTANT;
  }
}
