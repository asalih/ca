var requests = [
  {
    url: '/robots.txt',
    method: 'GET'
  }
];


function analyze(response) {

    if (response.URL.indexOf("robots.txt") > -1 && response.StatusCode == 200) {
      return Found(INFORMATION, "Robots.txt Detected")
    }
  }