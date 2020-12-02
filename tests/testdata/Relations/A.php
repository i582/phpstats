<?php

class RelatedClass {
  const CONSTANT = 0;

  public int $field = 10;

  public static function relatedMethod() {
    $tc = new TargetClass();
    $tc->targetMethod();
    $tc->some;
  }
}

class TargetClass extends RelatedClass {
  public int $some;

  public function targetMethod() {
    $tg = new RelatedClass();
    echo $tg->field;
    RelatedClass::relatedMethod();
    echo RelatedClass::CONSTANT;
  }
}
