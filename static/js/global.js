;(function ($) {

    $('[data-toggle="tooltip"]').tooltip();

    $.fn.open = function (title, values) {
      this.find(".modal-title").text(title);
      if(typeof values === "object"){
          for (var i in values){

          }
      }else{
          this.find(".modal-body").html(values)
      }
      return this;
    };

    var webHookModal = $("#webHookModal");
    var webHookCache = webHookModal.find(".modal-body").html();
    var serverModal = $("#serverModal");
    var serverCache = serverModal.find(".modal-body").html();

    $("#addWebHookBtn").on("click",function () {
        webHookModal.open("New WebHook",webHookCache).modal("show");
    });
    webHookModal.on("show.bs.modal",function () {
        $('[data-toggle="tooltip"]').tooltip();
    });

    $("#addServerBtn").on("click",function () {

        serverModal.open("New Server",serverCache).modal("show");
    });
    $("#serverForm").ajaxForm({
       beforeSubmit : function () {
            var serverNameEle = $("#serverName");

            if($.trim(serverNameEle.val()) === ""){
                serverNameEle.closest(".form-group").addClass("has-error");
                return false;
            }

            var serverIp = $("#serverIp");

            if($.trim(serverIp.val()) === ""){
                serverIp.closest(".form-group").addClass("has-error");
                return false;
            }
            var serverPort = $("#serverPort");
           if($.trim(serverPort.val()) === ""){
               serverPort.closest(".form-group").addClass("has-error");
               return false;
           }

           var serverKey = $("#serverKey");
           if($.trim(serverKey.val()) === ""){
               serverKey.closest(".form-group").addClass("has-error");
               return false;
           }

           $("#saveServerBtn").button("load");
       } ,
        success :function (res) {
            if (res.errcode === 0){
                if (serverModal.length > 0) {
                    var html = '<tr><td>'
                        + res.data.server_id + '</td><td>'
                        + res.data.name + '</td><td>'
                        + res.data.ip_address + '</td><td><span class="label label-success">正常</span></td><td>'
                        + res.data.account + '</td><td>'
                        + res.data.type + '</td><td>'
                        + res.data.time + '</td></tr>';

                    $("#serverTable>tbody").append(html);
                    serverModal.modal("hide");
                }else{
                    $("#errorMessage").css("color","green").text("success");
                }
            }else {
                $("#errorMessage").css("color","red").text(res.message);
            }
        },complete : function () {
            $("#saveServerBtn").button("reset");
        }
    });

    $("#webHookForm").ajaxForm({
        beforeSubmit :function () {

        },success : function (res) {

        },complate : function () {

        }
    })
})(jQuery);