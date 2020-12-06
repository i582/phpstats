<?php

class F {
  public I $IProp;

  public function FMethod(J $data): int {
    $g = new G; // ok
    $g->gProp; // ok
    H::$hProp; // ok
    $data->jProp; // ok
    $this->IProp->iProp; // ok
    return 0;
  }

  public function FMethod2(): int {
    return 0;
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
