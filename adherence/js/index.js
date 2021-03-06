
$.ajax({
    type: "POST",
    url: urlPath,
    success: function(data)
    {
        helpers.buildDropdown(
            jQuery.parseJSON(data),
            $('#dropdown'),
            'Select an option'
        );
    }
});


var helpers =
{
    buildDropdown: function(result, dropdown, emptyMessage)
    {
        // Remove current options
        dropdown.html('');
        
        // Add the empty option with the empty message
        dropdown.append('<option value="">' + emptyMessage + '</option>');

        // Check result isnt empty
        if(result != '')
        {
            // Loop through each of the results and append the option to the dropdown
            $.each(result, function(k, v) {
                dropdown.append('<option value="' + v.id + '">' + v.name + '</option>');
            });
        }
    }
}