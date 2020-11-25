<?php

class F {
  public I $IProp;

  public function FMethod(J $data) {
    $g = new G; // ok
    $g->gProp; // ok
    H::$hProp; // ok
    $data->jProp; // ok
    $this->IProp->iProp; // ok
  }
}

class G {
  public int $gProp = 0;
}

class H {
  public static int $hProp = 0;
}

class I {
  public int $iProp = 0;
}

class J {
  public int $jProp = 0;
}
