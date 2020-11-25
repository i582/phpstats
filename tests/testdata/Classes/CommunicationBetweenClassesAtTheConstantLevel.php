<?php

class K {
  const K_CONST = 0;

  public function KMethod() {
    N::N_CONST; // ok
  }
}

class L {
  const L_CONST = M::M_CONST;
}

class M {
  const M_CONST = K::K_CONST;
}

class N {
  const N_CONST = O::O_CONST;
}

class O {
  const O_CONST = L::L_CONST;
}
