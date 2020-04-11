var attacks = [
    {
      attack: '%2bprint(int)0xFFF9999-11205%3b%2f%2f',
    }
  ];
  
  function analyze(response) {

    if (response.Body.indexOf('26839803621') > -1) {
      return Found(CRITICAL, "Code Evaluation (PHP)")
    }
  }