;(function ($) {

    $.fn.open = function (title, values) {
      this.find(".modal-title").text(title);
      if(typeof values === "object"){
          for (var i in values){

          }
        }
      return this;
    };

    var webHookModal = $("#webHookModal");
    var serverModal = $("#serverModal");
    var serverCache = serverModal.find(".modal-body").html();

    $("#addWebHookBtn").on("click",function () {
        webHookModal.open("New WebHook","").modal("show");
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
       } ,
        success :function (res) {
            
        },error : function () {
            
        }
    });
})(jQuery);