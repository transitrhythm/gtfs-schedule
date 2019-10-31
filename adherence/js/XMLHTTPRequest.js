(function() {
    var httpRequest;
    function makeRequest(url) {
        if(window.XMLHttpRequest) { // Mozilla, Safari, ...
            httpRequest = new XMLHttpRequest();
        } else if (window.ActiveXObject) { // IE
            try {
                httpRequest = new ActiveXObject("Msxml2.XMPHTTP");
            }
            catch(e) {
                try {
                    httpRequest = new ActiveXObject("Microsoft.XMLHTTP");
                }
                catch(e) {}
            }
        }
        // call alertContents function after we receive server response
        httpRequest.onreadystatechange = function(){
            if (httpRequest.readyState === XMLHttpRequest.DONE) {
                if(httpRequest.status === 200) {
                    alert(httpRequest.responseText);
                } else {
                    alert('There was a problem with the request.');
                }
            }                
        }
        if (!httpRequest) {
            alert('Giving up :( Cannot create an XMLHTTP instance');
            return false;
        }
        // make the request
        httpRequest.open('GET', url);
        httpRequest.send();
    }
}