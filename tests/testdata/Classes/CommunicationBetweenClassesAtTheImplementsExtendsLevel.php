<?php

class P extends S implements Q {
  // ok
}

interface Q extends R {

}

interface R {

}

abstract class S extends T {

}

class T extends P implements Q {

}

