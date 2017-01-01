<div id="htpm-s2-objls-navbar">
  <ul id="htpm-s2-objls-dirnav" class="htpm-breadcrumb"></ul>
  <ul id="htpm-s2-objls-optools" class="htpm-node-nav htpm-nav-right">
    <li class="pure-button btapm-btn btapm-btn-primary">
      <a href="#" onclick="htpS2.ObjNew('file')">
        Upload New File
      </a>
    </li>
  </ul>
</div>

<div id="" class="htpm-div-light">
<table class="table table-hover">
  <thead>
    <tr>
      <th width="80px"></th>
      <th>Name</th>
      <th>Dir</th>
      <th>Size</th>
      <th></th>
      <th></th>
    </tr>
  </thead>
  <tbody id="htpm-s2-objls"></tbody>
</table>
</div>

<script id="htpm-s2-objls-dirnav-tpl" type="text/html">
{[~it.items :v]}
  <li><a href="#{[=v.path]}" onclick="htpS2.ObjList('{[=v.path]}')">{[=v.name]}</a></li>
{[~]}
</script>

<script id="htpm-s2-objls-tpl" type="text/html">  
  {[~it.items :v]}
    <tr id="obj{[=v._id]}">
      <td>
      {[ if (v.isdir) { ]}
        <span class="glyphicon glyphicon-folder-open" aria-hidden="true"></span>
      {[ } else if (v._isimg) { ]}
        <a href="{[=v.self_link]}" target="_blank"><img src="{[=v.self_link]}?ipn=i6040"></a>
      {[ } ]}
      </td>
      <td class="ts3-fontmono">
      {[ if (v.isdir) { ]}
        <a class="obj-item-dir" href="#objs" path="{[=v._abspath]}">{[=v.name]}</a>
      {[ } else { ]}
        <a class="obj-item-file" href="{[=v.self_link]}" target="_blank">{[=v.name]}</a>
      {[ } ]}
      </td>
      <td>
      {[?v.isdir]}YES{[?]}
      </td>
      <td>
      {[?!v.isdir]}
        {[=htpS2.UtilResourceSizeFormat(v.size)]}</td>
      {[?]}
      <td align="right">{[=l4i.TimeParseFormat(v.modtime, "Y-m-d H:i:s")]}</td>
      <td align="right">
      {[ if (!v.isdir) { ]}
        <a class="obj-item-del btn btn-default btn-xs" href="#obj-del" obj="{[=v._abspath]}">
          <span class="glyphicon glyphicon-cog" aria-hidden="true"></span> Delete
        </a>
      {[ } ]}
      </td>
    </tr>
  {[~]}
</script>

<!-- TPL : File New -->
<script id="htpm-s2-objnew-tpl" type="text/html"> 
<form id="{[=it.formid]}" action="#" onsubmit="htpS2.ObjNewSave('{[=it.formid]}');return false;">
<input type="hidden" name="type" value="{[=it.type]}">
<div class="form-group">
  <label>Folder Path</label>
  <input type="text" name="path" class="form-control" placeholder="Folder Path" value="{[=it.path]}">
</div>
<div class="form-group">
  <label>Select File</label>
  <input id="htpm-s2-objnew-files" type="file" name="file" class="form-control" placeholder="File Path" value="">
</div>
</form>
<div id="{[=it.formid]}-alert" class="alert alert-success" style="display:none"></div>
</script>

<!-- TPL : File Rename -->
<script id="htpm-s2-objrename-tpl" type="text/html"> 
<form id="{[=it.formid]}" action="#" onsubmit="htpS2.ObjRenameSave('{[=it.formid]}');return false;">
  <div class="input-prepend" style="margin-left:2px">
    <span class="add-on">
        <img src="{[=htpMgr.base]}-/img/folder_edit.png" class="h5c_icon">
    </span>
    <input type="text" name="pathset" value="{[=it.path]}" style="width:500px;">
    <input type="hidden" name="path" value="{[=it.path]}">
  </div>
</form>
</script>

<script type="text/javascript">
$("#htpm-s2-objls").on("click", ".obj-item-dir", function() {
    htpS2.ObjList($(this).attr("path"));
});
$("#htpm-s2-objls").on("click", ".obj-item-del", function() {
    var r = confirm("This file will be deleted, Confirm?");
    if (r == true) {
      htpS2.ObjDel($(this).attr("obj"));
    }
});
$("#htpS2-object-dirnav").on("click", ".obj-item-dir", function() {
    htpS2.ObjList($(this).attr("path"));
});
</script>