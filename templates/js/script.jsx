var Panel = React.createClass({
  render: function() {
    return (
      	<div className="row">
        	<div className="logo">
            <div className = "col-md-12" >
                <h1> Proxy-сервер </h1>
            </div>
          </div>
        </div>
      );
    }
  });

var Foot = React.createClass({
  render: function() {
    return (
    <footer>
       <div className="container">
          <div className="copy text-center">
             Copyright 2016 <a href='#'>Proxy site</a>
          </div>
       </div>
    </footer>
    )
  }
});

var Menu = React.createClass({
  getInitialState: function() {
    return {
      styleColl: null,
      collapse: ""
    };
  },
  collapseFunc: function() {
    if (this.state.styleColl) {
      this.setState({styleColl: null});
    } else {
      this.setState({styleColl: {display:'block'}});
    }
  },
  render: function() {
  return (
        <div className="col-md-3">
		  	<div className="sidebar content-box" style={{display:'block'}}>
                <ul className="nav">
                    <li className={(this.props.menuState === 1) ? "current" : ""} ><a href="#"><i className="glyphicon glyphicon-home"></i> Ваша статистика</a></li>
                    <li className={(this.props.menuState === 2) ? "current" : ""} onClick={this.props.change.bind(null, 2)}><a href="#"><i className="glyphicon glyphicon-calendar"></i> Calendar</a></li>
                    <li className={(this.props.menuState === 3) ? "current" : ""}><a href="#"><i className="glyphicon glyphicon-stats"></i> Statistics (Charts)</a></li>
                    <li className={(this.props.menuState === 4) ? "current" : ""}><a href="tables.html"><i className="glyphicon glyphicon-list"></i> Tables</a></li>
                    <li className={(this.props.menuState === 5) ? "current" : ""}><a href="buttons.html"><i className="glyphicon glyphicon-record"></i> Buttons</a></li>
                    <li className={(this.props.menuState === 6) ? "current" : ""}><a href="editors.html"><i className="glyphicon glyphicon-pencil"></i> Editors</a></li>
                    <li><a href="forms.html"><i className="glyphicon glyphicon-tasks"></i> Forms</a></li>
                    <li className={"submenu" + this.collapse} onClick={this.collapseFunc}>
                        <a href="#">
                            <i className="glyphicon glyphicon-list"></i> Pages
                            <span className="caret pull-right"></span>
                        </a>
                        <ul style={this.state.styleColl}>
                            <li><a href="login.html">Login</a></li>
                            <li><a href="signup.html">Signup</a></li>
                        </ul>
                    </li>
                </ul>
            </div>
		  </div>
    );
  }
});

var TableStat = React.createClass({
  render: function() {
      var sites = this.props.data.map(function(site) {
      return (
        <tr>
          <td>1</td>
          <td>{site.SiteName}</td>
          <td>{String(site.Traffic).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + ' байт'}</td>
        </tr>
      );
    });
    return (
      		<div className="row">
  				<div className="col-md-12">
  					<div className="content-box-large">
		  				<div className="panel-heading">
							<div className="panel-title">Sites</div>
							<div className="panel-options">
								<a href="#" data-rel="collapse"><i className="glyphicon glyphicon-refresh"></i></a>
								<a href="#" data-rel="reload"><i className="glyphicon glyphicon-cog"></i></a>
							</div>
						</div>
		  				<div className="panel-body">
		  					<table className="table">
				              <thead>
				                <tr>
				                  <th>#</th>
				                  <th>Site</th>
				                  <th>Traffic</th>
				                </tr>
				              </thead>
				              <tbody>
                      {sites}
				              </tbody>
				            </table>
		  				</div>
		  			</div>
  				</div>
  				</div>
      );
  }
});


var Main = React.createClass({
  // componentWillMount: function() {
  //   this.setState({idToken: null})
  // },
  getInitialState: function() {
    return {
      menuState: 1,
      data: []
    };
  },
  loadStats: function(id) {
    $.ajax({
      url: '../stat/1',
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data.Sites});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  componentDidMount: function() {
    this.loadStats();
  },
  changeMenuState: function(num) {
    this.setState({menuState: num})
  },
  render: function() {
    return (
      <div className="asdad">
      <div className = "container" >
        <Panel />
        <div className="row">
          <Menu menuState={this.state.menuState} change={this.changeMenuState}/>
            <div className="col-md-9">
            <TableStat data={this.state.data}/>
            </div>
        </div>
      </div>
      <Foot />
      </div>
      );
  }
});



    React.render( <Main /> ,
      document.getElementById('content')
    );
