var attacks = [
    {
      attack: '%27AND+1%3dcast(0x5f21403264696c656d6d61+as+varchar(8000))+or+%271%27%3d%27',
    }
  ];
  
  function analyze(response) {

    if (response.Body.indexOf('iNj3Ct3D') > -1) {
      console.log("Possible SQL Injection in " + response.URL)

      return true
    }

    return false
  }