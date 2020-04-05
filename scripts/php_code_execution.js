var attacks = [
    {
      attack: '+print(int)0xFFF9999-11205;//',
    }
  ];
  
  function analyze(response) {

    if (response.Body.indexOf('26839803621') > -1) {
      console.log("Code Evaluation (PHP) in " + response.URL)

      return true
    }

    return false
  }