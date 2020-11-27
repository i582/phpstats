<?php

class AA {
  public function AAMethod() {
    AB::ABMethod(); // ok
    $ac = new AC;
    $ac->ACMethod(); // ok
    AD::ADMethod(); // ok
    $ae = new AE;
    $ae->AEMethod(); // ok
  }
}

class AB {
  public static function ABMethod() {
    AD::ADMethod(); // ok
    $aa = new AA;
    $aa->AAMethod(); // ok
  }
}

class AC {
  public function ACMethod() {
    AB::ABMethod(); // ok
  }
}

class AD {
  public static function ADMethod() {
    $aa = new AA;
    $aa->AAMethod(); // ok
  }
}

class AE {
  public function AEMethod() {
    AD::ADMethod(); // ok
  }
}
