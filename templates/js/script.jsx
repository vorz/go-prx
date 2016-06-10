

var PlotDate = React.createClass({
  componentDidMount: function() {
        var legend = this.props.data.map(function(day){
          return [day.Day, day.Traffic/1000]
        });
        console.log(legend)
        $.jqplot('datePlot', [legend], {
              axes:{
                  xaxis:{
                      renderer:$.jqplot.DateAxisRenderer,
                      tickOptions:{formatString:'%d %a'}
                  }
              },
              series:[{lineWidth:4, markerOptions:{style:'square'}}]
        });
  },
  render: function() {
    return (
      <div className="row">
      <div className="col-md-10">
      <div class="affix">
        <div className="content-box-large">
          <div className="panel-heading">
          <div className="panel-title">График интернет активности за месяц</div>
        </div>
          <div className="panel-body">
            <div id="datePlot" style={{'margin-top':'20px', 'margin-left':'20px', 'width':'1000px', 'height':'300px'}}></div>
          </div>
        </div>
        </div>
      </div>
      </div>
    );
  }


})

var PlotPie = React.createClass({
  componentDidMount: function() {
        var legend = this.props.data.map(function(site){
          return [site.sitename, site.traffic]
        });
        $.jqplot('piePlot', [legend], {
            seriesDefaults:{
                shadow: true,
                renderer:$.jqplot.PieRenderer,
                rendererOptions:{
                    showDataLabels: true
                }
            },
            legend:{ show: true }
        });
  },
  render: function() {
    return (
      <div className="col-md-6">
      <div class="affix">
        <div className="content-box-large">
          <div className="panel-heading">
          <div className="panel-title">Диаграмма посещаемых сайтов</div>
        </div>
          <div className="panel-body">
            <div id="piePlot" style={{'margin-top':'20px', 'margin-left':'20px', 'width':'500px', 'height':'300px'}}></div>
          </div>
        </div>
        </div>
      </div>
    );
  }
});

var TableStat = React.createClass({
  getInitialState: function() {
    return {
      data: [],
      days: []
    };
  },
  loadStats: function(id) {

    $.ajax({
      url: '../stat/' + id ,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({data: data.Sites, days: data.Days});
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
        console.log(this.state.data)
        console.log(this.state.days)
        var PlotPieComponent
        if(this.state.data.length  != 0) {
          PlotPieComponent = <PlotPie data={this.state.data.slice(0,10)} />
        }
        var PlotDateComponent
        if(this.props.userInfo.UserID != -1 && this.state.days.length  != 0) {
          PlotDateComponent = <PlotDate data={this.state.days} />
        }
        var i = 0
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
          <div>
      		<div className="row">
  				<div className="col-md-6">
  					<div className="content-box-large">
		  				<div className="panel-heading">
							<div className="panel-title">Статистика трафика за месяц</div>
							<div className="panel-options">
								<a href="#" data-rel="collapse"><i className="glyphicon glyphicon-refresh"></i></a>
							</div>
						</div>
		  				<div className="panel-body">
		  					<table className="table">
				              <thead>
				                <tr>
				                  <th>#</th>
				                  <th>Сайт</th>
				                  <th>Трафик</th>
				                </tr>
				              </thead>
				              <tbody>
                      {sites.slice(0,10)}
				              </tbody>
				            </table>
		  				</div>
		  			</div>
  				</div>
          {PlotPieComponent}
  				</div>
          {PlotDateComponent}
          </div>
      );
  }
});

var Info = React.createClass({
  getInitialState: function() {
    return {
      withTable: false,
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
  render: function() {
    var TableComponent
    if(this.state.userInfo.length  != 0 && this.state.withTable) {
      TableComponent = <TableStat userInfo={this.state.userInfo} />
    }
    var TableLink
    if(!this.state.withTable) {
      TableLink = <h5><a href='#' onClick={function(){this.setState({withTable:true})}.bind(this)}>Посмотреть статистику</a></h5>
    }
    var traf = this.state.userInfo.Traffic ? "Трафик за текущий месяц: " + String(Math.floor(this.state.userInfo.Traffic/1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб  из ' + String(Math.floor(this.state.userInfo.Limit*1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб': "Пользователь не зарегистрирован."
    return (
      <div>
      <div className="row">
      <div className="col-md-6">
          <div className="content-box-header">
            <div className="panel-title"><strong>Информация о пользователе: {this.state.userInfo.IP}, {this.state.userInfo.Name}</strong></div>
          </div>
          <div className="content-box-large box-with-header">
            <p>{traf}</p>
            <br />
            {TableLink}
        </div>
        </div>
      </div>
      {TableComponent}
      </div>
    );
  }
});

//#################################################################################################

var Users = React.createClass({
  getInitialState: function() {
    return {
      users: []
    };
  },
  loadUserList: function() {
    $.ajax({
      url: '../users' ,
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({users: data.Users});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  componentDidMount: function() {
    this.loadUserList()
  },
  render: function() {
        var i = 0
        var users = this.state.users.map(function(user) {
        i++
        if(user.Name=='unknown') user.Name = 'нет'
        return (
          <tr>
            <td>{i}</td>
            <td><a href='#'>{user.IP}</a></td>
            <td>{user.Name}</td>
            <td>{String(Math.floor(user.Traffic/1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб'}</td>
            <td>{String(Math.floor(user.Limit*1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб'}</td>
            <td><button className="btn btn-danger btn-xs"><i className="glyphicon glyphicon-remove"></i> Удалить</button></td>
          </tr>
        );
    });
    return (
          <div>
      		<div className="row">
  				<div className="col-md-6">
  					<div className="content-box-large">
		  				<div className="panel-heading">
							<div className="panel-title">Пользователи</div>
						</div>
		  				<div className="panel-body">
		  					<table className="table">
				              <thead>
				                <tr>
				                  <th>#</th>
				                  <th>IP-адрес</th>
                          <th>DNS имя</th>
				                  <th>Трафик за месяц</th>
                          <th>Лимит</th>
				                  <th></th>
				                </tr>
				              </thead>
				              <tbody>
                      {users}
				              </tbody>
				            </table>
                    <br />
                    <button className="btn btn-primary btn-sm"><i className="glyphicon glyphicon-pencil"></i> Добавить пользователя</button>
		  				</div>
		  			</div>
  				</div>
  				</div>

          </div>
      );
  }
})

//############################################################################################################################

var OverallTable = React.createClass({
  render: function() {
        var lg = this.props.logs.map(function(stat) {
        return (
          <tr>
            <td>{stat.Date}</td>
            <td>{stat.SiteName}</td>
            <td>{stat.IP}</td>
            <td>{stat.Name}</td>
            <td>{String(Math.floor(stat.Traffic/1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб'}</td>
          </tr>
        );
    });
    return (
          <div>
      		<div className="row">
  				<div className="col-md-10">
  					<div className="content-box-large">
		  				<div className="panel-heading">
							<div className="panel-title">Последние соединения</div>
							<div className="panel-options">
								<a href="#" data-rel="collapse"><i className="glyphicon glyphicon-refresh"></i></a>
							</div>
						</div>
		  				<div className="panel-body">
		  					<table className="table">
				              <thead>
				                <tr>
				                  <th>Дата</th>
				                  <th>Сайт</th>
				                  <th>IP-адрес</th>
                          <th>Имя пользователя</th>
				                  <th>Количество байт</th>
				                </tr>
				              </thead>
				              <tbody>
                      {lg.slice(0,25)}
				              </tbody>
				            </table>
		  				</div>
		  			</div>
  				</div>
  				</div>
          </div>
      );
  }

})


var Overall = React.createClass({
  getInitialState: function() {
    return {
      stats: [],
    };
  },
  loadInfo: function() {
    $.ajax({
      url: '../overall',
      dataType: 'json',
      cache: false,
      success: function(data) {
        this.setState({stats: data});
      }.bind(this),
      error: function(xhr, status, err) {
        console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  },
  componentWillMount: function() {
    this.loadInfo();
  },
  render: function() {
    var oTable
    if (this.state.stats.length != 0) {
      oTable = <OverallTable logs={this.state.stats.Logs} />
    }
    return (
      <div>
      <div className="row">
      <div className="col-md-6">
          <div className="content-box-header">
            <div className="panel-title"><strong>Общая статистика сервера</strong></div>
          </div>
          <div className="content-box-large box-with-header">
          <ul class="list-unstyled">
            <li>Текущее количество TCP соединений: {this.state.stats.Conns}</li>
            <li>Количество зарегистрированных пользователей: {this.state.stats.UsersNum}</li>
            <li>Общий объем трафика за месяц: {String(Math.floor(this.state.stats.Traffic/1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб'}</li>
          </ul>
        </div>
        </div>
      </div>
      {oTable}
      </div>
    );
  }
})

//######################################################################################################################

var Sites = React.createClass({
  getInitialState: function() {
    return {
      sites: []
    };
  },
  // loadUserList: function() {
  //   $.ajax({
  //     url: '../users' ,
  //     dataType: 'json',
  //     cache: false,
  //     success: function(data) {
  //       this.setState({users: data.Users});
  //     }.bind(this),
  //     error: function(xhr, status, err) {
  //       console.error(this.props.url, status, err.toString());
  //     }.bind(this)
  //   });
  // },
  // componentDidMount: function() {
  //   this.loadUserList()
  // },
  render: function() {
    //     var i = 0
    //     var users = this.state.users.map(function(user) {
    //     i++
    //     if(user.Name=='unknown') user.Name = 'нет'
    //     return (
    //       <tr>
    //         <td>{i}</td>
    //         <td><a href='#'>{user.IP}</a></td>
    //         <td>{user.Name}</td>
    //         <td>{String(Math.floor(user.Traffic/1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб'}</td>
    //         <td>{String(Math.floor(user.Limit*1000)).replace(/(\d)(?=(\d{3})+([^\d]|$))/g, '$1 ') + '  Кб'}</td>
    //         <td><button className="btn btn-danger btn-xs"><i className="glyphicon glyphicon-remove"></i> Удалить</button></td>
    //       </tr>
    //     );
    // });
    var uI = {UserID: -1}
    return (
          <div>
      		<div className="row">
  				<div className="col-md-6">
  					<div className="content-box-large">
		  				<div className="panel-heading">
							<div className="panel-title">Сайты под фильтром</div>
						</div>
		  				<div className="panel-body">
		  					<table className="table">
				              <thead>
				                <tr>
				                  <th>#</th>
				                  <th>Название</th>
                          <th>Время</th>
				                  <th>Изменить</th>
				                  <th>Удалить</th>
				                </tr>
				              </thead>
				              <tbody>
                             <tr>
                               <td className="col-md-1">1</td>
                               <td className="col-md-4">sports.ru</td>
                               <td className="col-md-5">8:00 - 17:00</td>
                               <td className="col-md-1" className="text-center"><button className="btn btn-primary btn-xs"><i className="glyphicon glyphicon-pencil"></i> </button></td>
                               <td className="col-md-1" className="text-center"><button className="btn btn-danger btn-xs"><i className="glyphicon glyphicon-remove"></i></button></td>
                             </tr>
                             <tr>
                               <td className="col-md-1">2</td>
                               <td className="col-md-4">vk.com</td>
                               <td className="col-md-5">8:00 - 17:00</td>
                               <td className="col-md-1" className="text-center"><button className="btn btn-primary btn-xs"><i className="glyphicon glyphicon-pencil"></i> </button></td>
                               <td className="col-md-1" className="text-center"><button className="btn btn-danger btn-xs"><i className="glyphicon glyphicon-remove"></i></button></td>
                             </tr>
                             <tr>
                               <td className="col-md-1">3</td>
                               <td className="col-md-4">mail.ru</td>
                               <td className="col-md-5">8:00 - 17:00</td>
                               <td className="col-md-1" className="text-center"><button className="btn btn-primary btn-xs"><i className="glyphicon glyphicon-pencil"></i> </button></td>
                               <td className="col-md-1" className="text-center"><button className="btn btn-danger btn-xs"><i className="glyphicon glyphicon-remove"></i></button></td>
                             </tr>
				              </tbody>
				            </table>
                    <br />
                    <button className="btn btn-primary btn-sm"><i className="glyphicon glyphicon-plus"></i> Добавить сайт</button>
		  				</div>
		  			</div>
  				</div>
          <div className="col-md-6">
  					<div className="content-box-large">
		  				<div className="panel-heading">
							<div className="panel-title">Сайты вне учета</div>
						</div>
		  				<div className="panel-body">
		  					<table className="table">
				              <thead>
				                <tr>
				                  <th>#</th>
				                  <th>Название</th>
				                  <th>Изменить</th>
				                  <th>Удалить</th>
				                </tr>
				              </thead>
				              <tbody>
                             <tr>
                               <td className="col-md-1">1</td>
                               <td className="col-md-9">vlsu.ru</td>
                               <td className="col-md-1" className="text-center" ><button className="btn btn-primary btn-xs"><i className="glyphicon glyphicon-pencil"></i> </button></td>
                               <td className="col-md-1" className="text-center" ><button className="btn btn-danger btn-xs"><i className="glyphicon glyphicon-remove"></i> </button></td>
                             </tr>
				              </tbody>
				            </table>
                    <br />
                    <button className="btn btn-primary btn-sm"><i className="glyphicon glyphicon-plus"></i> Добавить сайт</button>
		  				</div>
		  			</div>
  				</div>
  				</div>
          <TableStat userInfo={uI} />
          </div>
      );
  }
})

var Main = React.createClass({
  getInitialState: function() {
    return {
      menuState: 1,
    };
  },
  changeMenuState: function(num, e) {
    e.preventDefault();
    this.setState({menuState: num})
  },
  render: function() {
    var content
    switch(this.state.menuState) {
    case 1:
      content =  <Info />;
      break;
    case 2:
      content = <Users />;
      break;
    case 3:
      content =  <Overall />;
      break;
    case 4:
      content =  <Sites />;
      break;
    }
    return (
      <div className="main">
      <div className = "container" >
        <Panel />
        <div className="content">
        <div className="row">
          <Menu menuState={this.state.menuState} change={this.changeMenuState} />
          <div className="col-md-10">
          {content}
          </div>
        </div>
        </div>
      </div>
      <Foot />
      </div>
      );
  }
});



var Panel = React.createClass({
  render: function() {
    return (
      <div class="header">
      	<div className="row">
        	<div className="logo">
          <div className = "col-md-1" >
              <img src="img/logo.png" />
          </div>
            <div className = "col-md-9" >
                <h1> Proxy-сервер </h1>
            </div>
	           <div class="col-md-2">
                <h3> Администратор <b className="caret" ></b></h3>
	           </div>
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
             Copyright 2016 <a href='#'>Go-prx</a>
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
                    <li className={(this.props.menuState === 1) ? "current" : ""} onClick={this.props.change.bind(null, 1)}><a href="#"><i className="glyphicon glyphicon-home"></i> Ваша статистика</a></li>
                    <li className={(this.props.menuState === 2) ? "current" : ""} onClick={this.props.change.bind(null, 2)}><a href="#"><i className="glyphicon glyphicon-calendar"></i> Пользователи</a></li>
                    <li className={(this.props.menuState === 3) ? "current" : ""} onClick={this.props.change.bind(null, 3)}><a href="#"><i className="glyphicon glyphicon-stats"></i> Общая статистика</a></li>
                    <li className={(this.props.menuState === 4) ? "current" : ""} onClick={this.props.change.bind(null, 4)}><a href="#"><i className="glyphicon glyphicon-stats"></i> Сайты </a></li>
                    <li className={(this.props.menuState === 5) ? "current" : ""} onClick={this.props.change.bind(null, 5)}><a href="#"><i className="glyphicon glyphicon-tasks"></i> Настройки</a></li>
                </ul>
            </div>
		  </div>

    );
  }
});
                    // <li className={(this.props.menuState === 6) ? "current" : ""}><a href="editors.html"><i className="glyphicon glyphicon-pencil"></i> Editors</a></li>
// <li><a href="forms.html"><i className="glyphicon glyphicon-tasks"></i> Forms</a></li>
//  <li className={"submenu" + this.collapse} onClick={this.collapseFunc}>
//     <a href="#">
//         <i className="glyphicon glyphicon-list"></i> Pages
//         <span className="caret pull-right"></span>
//     </a>
//     <ul style={this.state.styleColl}>
//         <li><a href="login.html">Login</a></li>
//         <li><a href="signup.html">Signup</a></li>
//     </ul>
// </li>


ReactDOM.render(<Main />,document.getElementById('content')
);
