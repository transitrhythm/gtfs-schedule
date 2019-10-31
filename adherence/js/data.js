function initialize_everything(callback) {
    // 2
    $.getJSON( "js/mysql_query_grants2.php", function(json){
        // 5
        // ...
        // NOW we have everything, so report back
        callback();
        }
    )
    // 3
}
