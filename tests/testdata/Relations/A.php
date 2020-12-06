<?php

class RelatedClass {
  const CONSTANT = 0;

  public int $field = 10;

  public static function relatedMethod(): float {
    $tc = new TargetClass();
    $tc->targetMethod();
    $tc->some;
    return 0.0;
  }
}


class TargetClass extends RelatedClass {
  public int $some;

  public function targetMethod(): int {
    $tg = new RelatedClass();
    echo $tg->field;
    RelatedClass::relatedMethod();
    echo RelatedClass::CONSTANT;
    return 0;
  }
}
