$("#formid").submit(function(event) {

  event.preventDefault();

  var $form = $(this),
    url = $form.attr('action');

  /* Send the data using post with element id name and name2*/
  var posting = $.post(url, {
    ip: $('#ip').val(),
    subnet: $('#subnet').val()
  });

  /* Alerts the results */
  posting.done(function(data) {
    alert('success');
    console.log(data);
    $('#result').html("<tr><th>Data Type</th><th>Value</th></tr>")
    $('#result').append(" \
      <tr> \
        <td>IP Address</td> \
        <td>" + data.IP + "</td>\
      </tr>");
    $('#result').append(" \
        <tr> \
          <td>Network Address</td> \
          <td>" + data.NetworkAddress + "</td>\
        </tr>");

    $('#result').append(" \
        <tr> \
          <td>Usable Host IP Range</td> \
          <td>" + data.Usable + "</td>\
        </tr>");

    $('#result').append(" \
        <tr> \
          <td>Broadcast Address</td> \
          <td>" + data.BroadcastAddress + "</td>\
        </tr>");

    $('#result').append(" \
        <tr> \
          <td>Number of Hosts</td> \
          <td>" + data.NumberOfHost + "</td>\
        </tr>");
    var usableHost = data.NumberOfHost > 2 ? data.NumberOfHost - 2 : 0
    $('#result').append(" \
        <tr> \
          <td>Number of Usable Hosts</td> \
          <td>" + usableHost+ "</td>\
        </tr>");

    $('#result').append(" \
          <tr> \
            <td>Subnet Mask</td> \
            <td>" + data.Subnet + "</td>\
          </tr>");

    $('#result').append(" \
          <tr> \
            <td>IP Class</td> \
            <td>" + "None" + "</td>\
          </tr>");

    $('#result').append(" \
          <tr> \
            <td>CIDR Notation</td> \
            <td>" +""+ "</td>\
          </tr>");

    $('#result').append(" \
          <tr> \
            <td>Short</td> \
            <td>" + data.IP + "</td>\
          </tr>");

    $('#result').append(" \
          <tr> \
            <td>Binary ID</td> \
            <td>" + data.BinID + "</td>\
          </tr>");

    $('#result').append(" \
            <tr> \
              <td>Integer ID</td> \
              <td>" + data.IntID + "</td>\
            </tr>");

    $('#result').append(" \
            <tr> \
              <td>Hexadecimal ID</td> \
              <td>" + data.HexID + "</td>\
            </tr>");

      $('#possible').html("<tr><th>Network Address</th><th>Usable Host Range</th><th>Broadcast Address</th></tr>")

      data.Possible.forEach(function(v){
        $('#possible').append("<tr><td>"+v.NetworkAddress+"</td><td>"+v.Usable+"</td><td>"+v.BroadcastAddress+"</td></tr>")
      })
  });

});
