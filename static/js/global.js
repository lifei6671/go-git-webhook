function loadCodeMirrorEditor(id) {
    window.CodeMirrorEditor = CodeMirror.fromTextArea(document.getElementById(id),{
        lineNumbers: true,
        mode: "text/x-sh",
        lineWrapping : true,
        matchBrackets: true,
        autoCloseBrackets: true
    });

    window.CodeMirrorEditor.on("change",function () {
       var content = window.CodeMirrorEditor.getValue();
       $("#" + id).val(content);
    });
}

(function ($) {

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
    webHookModal.on("shown.bs.modal",function () {
        if(!window.CodeMirrorEditor){
            window.loadCodeMirrorEditor("shellScript");
        }
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

                    $("#serverTable>tbody").prepend(res.view);
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
            var isValid = true;
            var repositoryNameEle = $("#repositoryName");
            if($.trim(repositoryNameEle.val()) == ""){
                repositoryNameEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }
            var repositoryBranchEle = $("#repositoryBranch");

            if($.trim(repositoryBranchEle.val()) === ""){
                repositoryBranchEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }

            var serverTagEle = $("#serverTag");

            if($.trim(serverTagEle.val()) === ""){
                serverTagEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }
            var shellScriptEle = $("#shellScript");

            if($.trim(shellScriptEle.val()) === ""){
                shellScriptEle.closest(".form-group").addClass("has-error");
                isValid = false;
            }

            if(!isValid){
                return false;
            }
            $("#saveWebHookBtn").button("load");
        },success : function (res) {
            if (res.errcode === 0){
                if (webHookModal.length > 0) {
                    $("#webHookTable>tbody").prepend(res.view);
                    webHookModal.modal("hide");
                }else{
                    $("#errorMessage").css("color","green").text("success");
                }
            }else {
                $("#errorMessage").css("color","red").text(res.message);
            }
        },complete : function () {
            $("#saveWebHookBtn").button("reset");
        }
    });
    
    
})(jQuery);