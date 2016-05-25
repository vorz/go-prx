



var Plot = React.createClass({
  componentDidMount: function() {
        $.jqplot('piePlot', [[['a',34000],['b',75000],['c',2000]]], {
            seriesDefaults:{
                shadow: true,
                renderer:$.jqplot.PieRenderer,
                rendererOptions:{
                    sliceMargin: 4,
                    // rotate the starting position of the pie around to 12 o'clock.
                    startAngle: -90
                }
            },
            legend:{ show: true }
        });
  },
  render: function() {
    return (
      <div className="col-md-5">
        <div className="content-box-large">
          <div className="panel-heading">
          <div className="panel-title">Pie Plot</div>
        </div>
          <div className="panel-body">
            <div id="piePlot" style={{'margin-top':'20px', 'margin-left':'20px', 'width':'400px', 'height':'300px'}}></div>
          </div>
        </div>
      </div>
    );
  }
});

var TableStat = React.createClass({
  getInitialState: function() {
    return {
      data: []
    };
  },
  loadStats: function(id) {
    $.ajax({
      url: '../stat/' + id ,
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
    this.loadStats(this.props.userInfo.UserID)
  },
  render: function() {
      var i = 0
          console.log(this.props.userInfo)
      var sites = this.state.data.map(function(site) {
        i++
      return (
        <tr>
          <td>{i}</td>
          <td>{site.sitename}</td>
          <td>{String(Math.floor(site.traffic/1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб'}</td>
        </tr>
      );
    });
    return (
      		<div className="row">
  				<div className="col-md-7">
  					<div className="content-box-large">
		  				<div className="panel-heading">
							<div className="panel-title">Sites</div>
							<div className="panel-options">
								<a href="#" data-rel="collapse"><i className="glyphicon glyphicon-refresh"></i></a>
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
                      {sites.slice(0,10)}
				              </tbody>
				            </table>
		  				</div>
		  			</div>
  				</div>
          <Plot />
  				</div>
      );
  }
});

var Info = React.createClass({
  render: function() {
    var traf = this.props.userInfo.Traffic ? "Трафик за текущий месяц: " + this.props.userInfo.Traffic : "Пользователь не зарегистрирован."
    return (
      <div className="row">
      <div className="col-md-6">
          <div className="content-box-header">
            <div className="panel-title"><strong>Информация о пользователе: {this.props.userInfo.IP}</strong></div>
          </div>
          <div className="content-box-large box-with-header">
            {traf}
          <br /><br />
        </div>
        </div>
      </div>
    );
  }
})

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
        <div className="col-md-2">
		  	<div className="sidebar content-box" style={{display:'block'}}>
                <ul className="nav">
                    <li className={(this.props.menuState === 1) ? "current" : ""} ><a href="/table"><i className="glyphicon glyphicon-home"></i> Ваша статистика</a></li>
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

var Main = React.createClass({
  // componentWillMount: function() {
  //   this.setState({idToken: null})
  // },
  getInitialState: function() {
    return {
      menuState: 1,
      userInfo: []
    };
  },
  loadInfo: function() {
    $.ajax({
      url: '../info',
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({userInfo: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  componentWillMount: function() {
    this.loadInfo();
  },
  changeMenuState: function(num) {
    this.setState({menuState: num})
  },
  render: function() {
    var TableComponent
    if(this.state.userInfo.length  != 0) {
      TableComponent = <TableStat userInfo={this.state.userInfo} />
    }
    return (
      <div className="main">
      <div className = "container" >
        <Panel />

        <div className="row">
          <Menu menuState={this.state.menuState} change={this.changeMenuState} />
          <div className="col-md-10">
          <Info userInfo={this.state.userInfo} />
          {TableComponent}
          </div>
        </div>
      </div>
      <Foot />
      </div>
      );
  }
});



ReactDOM.render(<Main />,document.getElementById('content')
);
