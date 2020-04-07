var attacks = [
    {
      attack: '+print(int)0xFFF9999-11205;//',
    }
  ];
  
  function analyze(response) {

    if (response.Body.indexOf('26839803621') > -1) {
      return Found(CRITICAL, "Code Evaluation (PHP)")
    }
  }