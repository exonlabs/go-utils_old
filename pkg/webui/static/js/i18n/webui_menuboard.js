/*
  :copyright: 2019-2023, ExonLabs. All rights reserved.
  :license: BSD, see LICENSE for more details.
*/
var WebUI = function($, ui) {

  ui.scrolltop = function(interval) {
    $("#board-page").animate({scrollTop:0},(interval)?interval:300);
  };

  ui.board_menu = {
    show: function() {
      $("body").addClass("MenuToggled");
    },
    hide: function() {
      $("body").removeClass("MenuToggled");
    },
    toggle: function() {
      $("body").toggleClass("MenuToggled");
    },
    show_submenu: function(id) {
      var m=$('#board-menubody #submenu_'+id);
      if(m.attr('aria-expanded')!='true') m.click();
    },
    hide_submenu: function(id) {
      var m=$('#board-menubody #submenu_'+id);
      if(m.attr('aria-expanded')=='true') m.click();
    },
    toggle_submenu: function(id) {
      $('#board-menubody #submenu_'+id).click();
    }
  };

  ui.board_content = {
    old_hash: null,
    load_neglect: false,
    update: function(data) {
      $("#pagebody-contents").html(data);
      ui.scrolltop(0);
    },
    error: function(message) {
      ui.board_content.update(
        '<div class="p-3"><div class="alert alert-danger text-left">' +
        '<i class="fa fas fa-ta fa-exclamation-circle"></i> ' + message + '</div></div>');
    },
    warn: function(message) {
      ui.board_content.update(
        '<div class="p-3"><div class="alert alert-warning text-left">' +
        '<i class="fa fas fa-ta fa-exclamation-circle"></i> ' + message + '</div></div>');
    },
    info: function(message) {
      ui.board_content.update(
        '<div class="p-3"><div class="alert alert-info text-left">' +
        '<i class="fa fas fa-ta fa-info-circle"></i> ' + message + '</div></div>');
    },
    success: function(message) {
      ui.board_content.update(
        '<div class="p-3"><div class="alert alert-success text-left">' +
        '<i class="fa fas fa-ta fa-check-circle"></i> ' + message + '</div></div>');
    },
    load: function(verb, url, params) {
      if(ui.board_content.load_neglect) {
        ui.board_content.load_neglect = false;
        return null;
      };
      ui.loader.load(verb, url, params,
        function(result) {
          ui.board_content.old_hash = window.location.hash;
          if(result.redirect) ui.redirect(result.redirect, result.blank);
          else {
            if(result.doctitle) ui.doctitle.update(result.doctitle);
            if(result.notifications) ui.notify.load(result.notifications);
            if(result.payload !== undefined) ui.board_content.update(result.payload);
          };
        },
        function(error) {
          if(error !== null) ui.notify.error(error,true,false);
          if(ui.board_content.old_hash) {
            ui.board_content.load_neglect = true;
            window.location.hash = ui.board_content.old_hash;
          };
        }
      );
    }
  };

  $(document).ready(function() {
    $(window)
      .on("resize", function() {
        if(window.innerWidth < 992) ui.board_menu.hide();
      })
      .on("hashchange", function(e) {
        e.preventDefault();
        if(window.innerWidth < 992) ui.board_menu.hide();
        ui.board_content.load("GET", window.location.hash, null);
      });

    $("#board-menubody>ul.metismenu").metisMenu();

    $("body")
      .on("click", "a.pagelink", function(e) {
        e.preventDefault();
        ui.redirect($(this).attr("href"));
      });

    $("#board-menutoggle>a")
      .bind("click", function(e) {
        e.preventDefault();
        ui.board_menu.toggle();
      });

    $("#board-backdrop")
      .bind('click', function(e) {
        e.preventDefault();
        ui.board_menu.hide();
      });

    $('#board-wrapper').show();
    setTimeout(function() {
      if(window.location.hash.length <= 1) {
        var loc = $('#board-menubody a.pagelink').attr("href");
        if(loc !== undefined) window.location.hash = loc;
      }
      else $(window).trigger("hashchange");
    }, 200);
  });

  return ui;
}(jQuery, WebUI || {});
