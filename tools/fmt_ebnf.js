function processTextNodes(node, ruleID, pattern) {
    if (node.nodeType === Node.TEXT_NODE) {
      var replaced = node.textContent.replace(pattern, function(match) {
        if (match === ruleID) {
          // Self-reference becomes a clickable span.
          return `<span class="rule-name" data-rule="${ruleID}" onclick="showDeps(event)">${match}</span>`;
        } else {
          // Other rule references become anchor links.
          return `<a href="#${match}" data-rule="${match}">${match}</a>`;
        }
      });
      var span = document.createElement("span");
      span.innerHTML = replaced;
      node.parentNode.replaceChild(span, node);
    } else if (node.nodeType === Node.ELEMENT_NODE && !node.classList.contains("rule-name")) {
      Array.from(node.childNodes).forEach(function(child) {
        processTextNodes(child, ruleID, pattern);
      });
    }
  }
  
  document.addEventListener("DOMContentLoaded", function() {
    var ruleNames = Object.keys(dependents).sort();
    var pattern = new RegExp("\\b(" + ruleNames.join("|") + ")\\b", "g");
    document.querySelectorAll("div.rule > code").forEach(function(codeElem) {
      var ruleID = codeElem.parentElement.id;
      processTextNodes(codeElem, ruleID, pattern);
    });
    
    // Attach tooltip events to all anchor links.
    document.querySelectorAll("div.rule a[data-rule]").forEach(function(linkElem) {
      linkElem.addEventListener("mouseenter", showPreview);
      linkElem.addEventListener("mouseleave", hidePreview);
    });
  });
  
  function showDeps(e) {
    var ruleName = e.target.dataset.rule;
    var deps = dependents[ruleName] || [];
    var message = "";
    if (deps.length === 0) {
      message = `<strong>${ruleName}</strong> is not referenced by any other rule.`;
    } else {
      message = `<strong>Rules that depend on ${ruleName}:</strong><br><ul>`;
      deps.forEach(function(dep) {
        message += `<li><a href="#${dep}">${dep}</a></li>`;
      });
      message += `</ul>`;
    }
    var popup = document.getElementById("dependencyPopup");
    popup.innerHTML = message;
    var rect = e.target.getBoundingClientRect();
    popup.style.display = "block";
    popup.style.top = (window.scrollY + rect.bottom + 5) + "px";
    popup.style.left = (window.scrollX + rect.left) + "px";
    e.stopPropagation();
  }
  
  document.addEventListener("click", function(e) {
    var popup = document.getElementById("dependencyPopup");
    if (popup.style.display === "block") {
      var isPopup = popup.contains(e.target);
      var isRuleName = e.target.classList && e.target.classList.contains("rule-name");
      if (!isPopup && !isRuleName) {
        popup.style.display = "none";
      }
    }
  });
  
  function showPreview(e) {
    var ruleName = e.target.dataset.rule;
    if (!ruleName) return;
    var content = ruleContents[ruleName] || "(No content)";
    var popup = document.getElementById("previewPopup");
    popup.innerHTML = `<code>${content}</code>`;
    var rect = e.target.getBoundingClientRect();
    popup.style.display = "block";
    popup.style.top = (window.scrollY + rect.bottom + 5) + "px";
    popup.style.left = (window.scrollX + rect.left) + "px";
  }
  
  function hidePreview(_) {
    var popup = document.getElementById("previewPopup");
    popup.style.display = "none";
  }
