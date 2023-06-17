/*
  :copyright: 2019-2023, ExonLabs. All rights reserved.
  :license: BSD, see LICENSE for more details.
*/
var WebUI = function($, ui) {

  $.ajaxSetup({cache:true});

  ui.format_url = function(url) {
    return "/" + url.replace(/^[\/#]/, "");
  };

  ui.doctitle = {
    base: document.title,
    update: function(title) {
      if(title !== undefined && title.length > 0) {
        title = title[0].toUpperCase() + title.slice(1);
        document.title = title + " | " + ui.doctitle.base;
      };
    },
    reset: function() {
      document.title = ui.doctitle.base;
    }
  };

  ui.pagelock = {
    show: function(html, styles) {
      ui.pagelock.hide();
      $("body").append('<div id="_UiPageLock" class="container-fluid overflow-auto scroll page-lock '+(styles?styles:'')+'">'+(html?html:'')+'</div>');
    },
    loading: function(html) {
      ui.pagelock.show('<div class="row"><div class="col cancel"><a id="_UiPageLock_btnCancel">&times;</a></div></div>'+(html?html:''),'page-loading');
      $("#_UiPageLock_btnCancel").on('click',function(){ui.pagelock.hide()});
      return $("#_UiPageLock_btnCancel");
    },
    progress: function(percent) {
      if(percent) $("#_UiPageLock .progress-bar").css("width", percent+"%");
      else return ui.pagelock.loading('<div class="row h-100 align-items-center justify-content-center"><div class="col-5"><div class="progress"><div class="progress-bar progress-bar-striped progress-bar-animated"></div></div></div></div>');
    },
    modal: function(title, contents, footer, styles) {
      ui.pagelock.show('<div class="modal-dialog modal-dialog-centered '+(styles?styles:'')+'"><div class="modal-content"><div class="modal-header pt-2 pb-1">'+(title?title:'')+'<button class="close" onclick="WebUI.pagelock.hide()">&times;</button></div><div class="modal-body scroll">'+(contents?contents:'')+(footer?'</div><div class="modal-footer p-1">'+footer+'</div></div></div>':''));
    },
    hide: function() {
      if($("#_UiPageLock").length) $("#_UiPageLock").remove();
    }
  };

  ui.notify = {
    stack: {"dir1":"down", "dir2":"left", "push":"top",
            "firstpos1":14, "firstpos2":10, "spacing1":7, "spacing2":7},
    clear: PNotify.removeAll,
    show: function(category, message, unique, sticky) {
      if(unique) PNotify.removeAll();
      if(window.innerWidth < 576) {
        ui.notify.stack.firstpos1 = 0; ui.notify.stack.firstpos2 = 0;
      } else {
        ui.notify.stack.firstpos1 = 14; ui.notify.stack.firstpos2 = 10;
      };
      if($("html").attr("dir") == "rtl" || $("body").attr("dir") == "rtl")
        ui.notify.stack.dir2 = "right";
      else ui.notify.stack.dir2 = "left";
      var opt = {
        styling: "fontawesome", icon: false, hide: (sticky?false:true),
        animate_speed: "fast", buttons: {sticker:false, closer:true},
        addclass: "stack-custom", stack: ui.notify.stack
      };
      if(category == "error") {
        opt.type = "error";
        opt.text = '<i class="fa fas fa-ta fa-exclamation-circle"></i>' + message;
      } else if(category == "warn") {
        opt.type = "notice";
        opt.text = '<i class="fa fas fa-ta fa-exclamation-circle"></i>' + message;
      } else if(category == "success") {
        opt.type = "success";
        opt.text = '<i class="fa fas fa-ta fa-check-circle"></i>' + message;
      } else {
        opt.type = "info";
        opt.text = '<i class="fa fas fa-ta fa-info-circle"></i>' + message;
      };
      var n = new PNotify(opt);
      n.get().click(function() {n.remove()});
    },
    error: function(message, unique, sticky) {
      ui.notify.show('error', message, unique, sticky);
    },
    warn: function(message, unique, sticky) {
      ui.notify.show('warn', message, unique, sticky);
    },
    info: function(message, unique, sticky) {
      ui.notify.show('info', message, unique, sticky);
    },
    success: function(message, unique, sticky) {
      ui.notify.show('success', message, unique, sticky);
    },
    load: function(notifications) {
      for(var i=0; i<notifications.length; i++)
        ui.notify.show(notifications[i][0],notifications[i][1],notifications[i][2],notifications[i][3]);
    }
  };

  ui.scrolltop = function(interval) {
    $("body").animate({scrollTop:0},(interval)?interval:300);
  };

  ui.redirect = function(url, blank) {
    if(url !== undefined && url.length > 0) {
      if(url[0] == '#') {
        if(url == window.location.hash) $(window).trigger("hashchange");
        else window.location.hash = url;
      }
      else if(blank) window.open(url);
      else if(url == window.location) location.reload();
      else window.location = url;
    };
  };

  ui.request = function(verb, url, params, fSuccess, fError, fComplete, fXhr) {
    return $.ajax({
      url: ui.format_url(url), type: verb, data: params?params:{},
      contentType:(params instanceof FormData)?false:'application/x-www-form-urlencoded; charset=UTF-8',
      processData:!(params instanceof FormData),
      success: function(result, status, xhr) {
        ui.request.success(result, fSuccess);
      },
      error: function(xhr, status, error) {
        if(error == 'abort') error = null;
        else if(!xhr.status) error = $.i18n._("no connection");
        else if(!error) error = $.i18n._("request failed");
        ui.request.error(error, fError);
      },
      complete: function(xhr, status) {
        ui.request.complete(status, fComplete);
      },
      xhr: function() {
        if(typeof fXhr === "function") return fXhr();
        else return (new window.XMLHttpRequest());
      }
    });
  };
  ui.request.success = function(result, fSuccess) {
    if(typeof fSuccess === "function") fSuccess(result);
    else {
      if(result.notifications) ui.notify.load(result.notifications);
      if(result.redirect) ui.redirect(result.redirect, result.blank);
    };
  };
  ui.request.error = function(error, fError) {
    if(typeof fError === "function") fError(error);
    else if(error !== null) ui.notify.error(error,true,false);
  };
  ui.request.complete = function(status, fComplete) {
    if(typeof fComplete === "function") fComplete(status);
  };

  ui.loader = {
    req_xhr: null,
    lock_timer: null,
    progress_timer: null,
    progress_xhr: null,
    load: function(verb, url, params, fSuccess, fError, fComplete, timeout) {
      ui.loader.cancel();
      ui.loader.lock_timer = setTimeout(function() {
        ui.pagelock.loading().off("click").on("click", function(e) {ui.loader.cancel()});
      }, (timeout)?timeout:500);
      ui.loader.req_xhr = ui.request(verb, url, params,
        function(result) {
          ui.loader.reset();
          ui.pagelock.hide();
          ui.request.success(result, fSuccess);
        },
        function(error) {
          ui.loader.reset();
          ui.pagelock.hide();
          ui.request.error(error, fError);
        },
        fComplete
      );
    },
    progress: function(verb, url, params, fSuccess, fError, fComplete, timeout, interval) {
      ui.loader.cancel();
      ui.loader.lock_timer = setTimeout(function() {
        ui.pagelock.progress().off("click").on("click", function(e) {ui.loader.cancel()});
        ui.loader.progress_timer = setInterval(function() {
          if(ui.loader.progress_xhr === null) {
            ui.loader.progress_xhr = ui.request(verb, url, {get_progress:1},
              function(r){ui.pagelock.progress(r.payload)}, function(e){}, function(s){ui.loader.progress_xhr=null});
          };
        }, (interval)?interval:5000);
      }, (timeout)?timeout:500);
      ui.loader.req_xhr = ui.request(verb, url, params,
        function(result) {
          ui.loader.reset();
          ui.pagelock.progress(100);
          setTimeout(function() {
            ui.pagelock.hide();
            ui.request.success(result, fSuccess);
          }, 500);
        },
        function(error) {
          ui.loader.reset();
          ui.pagelock.hide();
          ui.request.error(error, fError);
        },
        fComplete
      );
    },
    formsubmit: function(form, fSuccess, fError, fComplete, timeout) {
      ui.loader.cancel();
      var uiloader = ui.pagelock.progress;
      ui.loader.lock_timer = setTimeout(function() {
        uiloader().off("click").on("click", function(e) {ui.loader.cancel()});
      }, (timeout)?timeout:500);
      ui.loader.req_xhr = ui.request(
        form.attr('method'), form.attr('action'), (new FormData(form[0])),
        function(result) {
          ui.loader.reset();
          ui.pagelock.hide();
          ui.request.success(result, fSuccess);
        },
        function(error) {
          ui.loader.reset();
          ui.pagelock.hide();
          ui.request.error(error, fError);
        },
        fComplete,
        function() {
          var xhr = new window.XMLHttpRequest();
          xhr.upload.addEventListener("progress", function(evt) {
            if (evt.lengthComputable) {
              var percent = parseInt(evt.loaded / evt.total * 100);
              if(percent<100) {ui.pagelock.progress(percent);return null};
              ui.pagelock.progress(100);
            };
            uiloader = ui.pagelock.loading;
          }, false);
          return xhr;
        }
      );
    },
    cancel: function() {
      if(ui.loader.req_xhr) ui.loader.req_xhr.abort();
    },
    reset: function() {
      if(ui.loader.lock_timer) clearTimeout(ui.loader.lock_timer);
      if(ui.loader.progress_timer) clearTimeout(ui.loader.progress_timer);
      ui.loader.req_xhr = null;
      ui.loader.lock_timer = null;
      ui.loader.progress_timer = null;
    }
  };

  ui.static_local_url_prefix = "/static/vendor";

  ui.loadCss = function(url, altUrl) {
    if(!$("head link[href='"+url+"']").length) {
      $('head').append('<link rel="stylesheet" type="text/css" href="'+url+'" '+(altUrl?'onerror="this.onerror=null;this.href=\''+altUrl+'\'"':'')+'>');
    };
  };
  ui.loadCssCdn = function(cdnUrl, resourcePath) {
    if(!cdnUrl) ui.loadCss(ui.static_local_url_prefix+'/'+resourcePath);
    else ui.loadCss(cdnUrl+'/'+resourcePath,ui.static_local_url_prefix+'/'+resourcePath);
  };

  ui.loadScript = function(url, fSuccess, fError) {
    if($("body script[_src='"+url+"']").length) {
      if(typeof fSuccess === "function") return fSuccess()};
    $.getScript(url)
      .done(function(script, status){
        $('body').append('<script type="text/javascript" _src="'+url+'"></script>');
        if(typeof fSuccess === "function") fSuccess();
      })
      .fail(function(xhr, status, error) {
        if(typeof fError === "function") fError();
        else ui.notify.warn($.i18n._('Failed to load all page contents !!') + '<br>' + $.i18n._('Please reload page and try again.'),null,true);
      });
  };
  ui.loadScriptCdn = function(cdnUrl, resourcePath, fSuccess, fError) {
    if(!cdnUrl) return ui.loadScript(ui.static_local_url_prefix+'/'+resourcePath, fSuccess, fError);
    ui.loadScript(cdnUrl+'/'+resourcePath, fSuccess, function(){
      ui.loadScript(ui.static_local_url_prefix+'/'+resourcePath, fSuccess, fError);
    });
  };

  $(document).ready(function() {
    var lang = $('html').attr("lang");
    if(lang && typeof(webui_i18n) != "undefined") $.i18n.load(webui_i18n);
  });

  return ui;
}(jQuery, WebUI || {});
