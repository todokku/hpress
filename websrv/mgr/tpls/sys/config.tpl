<style>
.htp-sys-table {
  font-size: 12px;
}
.htp-sys-table td {
  padding: 5px !important;
}
.htp-sys-table tr.line {
  border-top: 1px solid #ccc;
}
</style>


<div class="panel panel-default">
  <div class="panel-heading">System Config</div>
  <div id="htpm-sys-configset" class="panel-body">

    <div id="htpm-sys-configset-alert"></div>

    <table width="100%" class="table table-striped">
      <thead>
        <tr>
          <th>Key</th>
          <th>Value</th>
          <th>Comment</th>
        </tr>
      </thead>
      <tbody>
      {[~it.items :v]}
      <tr>
        <td width="20%">{[=v.key]}</td>
        <td width="40%">
          {[ if (v.type !== undefined && v.type == "text") { ]}
          <textarea class="form-control htpm-sys-config-item" name="{[=v.key]}" rows="3">{[=v.value]}</textarea>
          {[ } else { ]}
          <input type="text" class="form-control htpm-sys-config-item" name="{[=v.key]}" value="{[=v.value]}">
          {[ } ]}
        </td>
        <td>{[=v.comment]}</td>
      </tr>
      {[~]}
      </tbody>
    </table>

    <button class="pure-button btapm-btn btapm-btn-primary" onclick="htpSys.ConfigSetCommit()">Save</button>
  </div>
</div>
