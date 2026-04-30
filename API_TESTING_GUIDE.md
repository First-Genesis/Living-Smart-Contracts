# Living Smart Contracts API Testing Guide

## 🚀 Quick Start

The Living Smart Contracts API is now running locally with full Swagger documentation and interactive testing capabilities.

### Server Status
- **Status**: ✅ Running
- **Port**: 8080
- **Swagger UI**: http://localhost:8080/swagger/
- **Health Check**: http://localhost:8080/health

## 📊 Available Endpoints

### Health & Information
- `GET /health` - Health check
- `GET /api/info` - API information and available endpoints

### Contract Management
- `GET /api/contracts` - List all contracts
- `POST /api/contracts/deploy` - Deploy a new contract
- `POST /api/contracts/execute` - Execute contract function
- `GET /api/contracts/{address}` - Get contract details

### Evolution & Learning
- `POST /api/evolution/trigger` - Trigger contract evolution

### Collaboration
- `POST /api/collaborations/propose` - Propose contract collaboration

## 🧪 Testing Methods

### 1. Interactive Swagger UI
Visit **http://localhost:8080/swagger/** for interactive API testing with:
- Complete API documentation
- Request/response examples
- Try-it-out functionality
- Schema validation

### 2. Automated Test Script
Run the comprehensive test script:
```bash
./test_api.sh
```

### 3. Manual cURL Commands

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Deploy Contract
```bash
curl -X POST http://localhost:8080/api/contracts/deploy \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my_living_contract",
    "type": "living",
    "source_code": "contract MyContract { function hello() public pure returns (string) { return \"Hello World\"; } }",
    "owner": "0x1234567890123456789012345678901234567890",
    "init_params": {"initial_value": 100},
    "time_aware": true,
    "history_depth": 1000
  }'
```

#### Execute Contract
```bash
curl -X POST http://localhost:8080/api/contracts/execute \
  -H "Content-Type: application/json" \
  -d '{
    "contract_address": "0xabcdef1234567890abcdef1234567890abcdef12",
    "function": "hello",
    "parameters": {},
    "caller": "0x1234567890123456789012345678901234567890",
    "gas_limit": 100000
  }'
```

#### Trigger Evolution
```bash
curl -X POST http://localhost:8080/api/evolution/trigger \
  -H "Content-Type: application/json" \
  -d '{
    "evolution_type": "optimization",
    "parameters": {
      "target_metric": "gas_efficiency",
      "mutation_rate": 0.05,
      "selection_pressure": 0.8
    }
  }'
```

## 📋 API Features Demonstrated

### ✅ Core Features
- **RESTful API Design** - Standard HTTP methods and status codes
- **JSON Request/Response** - Structured data exchange
- **CORS Support** - Cross-origin resource sharing enabled
- **Error Handling** - Proper error responses with codes and messages

### ✅ Living Smart Contracts Features
- **Contract Deployment** - Deploy new living contracts
- **Contract Execution** - Execute functions with parameters
- **Evolution Triggers** - Initiate evolutionary processes
- **Collaboration Proposals** - Inter-contract partnerships
- **Performance Metrics** - Success rates, gas usage, adaptation scores
- **DNA & Genetics** - Contract genetic information and fitness
- **Memory Systems** - Short-term and long-term contract memory

### ✅ Documentation & Testing
- **Swagger/OpenAPI 3.0** - Complete API documentation
- **Interactive UI** - Test endpoints directly in browser
- **Request Validation** - Schema validation for all inputs
- **Response Examples** - Sample responses for all endpoints

## 🔧 Development Features

### Mock Data
The API currently returns realistic mock data including:
- Sample living contracts with DNA and performance metrics
- Evolutionary genetics with genes, mutations, and fitness scores
- Collaboration networks and partnership data
- Memory patterns and learning experiences

### Extensibility
The API is designed for easy extension:
- Modular handler functions
- Structured data models
- Comprehensive error handling
- CORS-enabled for frontend integration

## 🌐 Browser Testing

1. **Swagger UI**: http://localhost:8080/swagger/
   - Interactive API documentation
   - Test all endpoints directly
   - View request/response schemas

2. **Health Check**: http://localhost:8080/health
   - Quick server status verification

3. **API Info**: http://localhost:8080/api/info
   - Service information and available endpoints

## 📈 Performance Metrics

The API demonstrates Living Smart Contracts concepts including:
- **Execution Count**: Number of function calls
- **Success Rate**: Percentage of successful executions
- **Gas Efficiency**: Average gas consumption
- **Adaptation Score**: Learning and optimization metrics
- **Fitness Score**: Evolutionary fitness rating
- **Generation Tracking**: Contract evolution generations

## 🚀 Next Steps

1. **Frontend Integration**: Use the API with web applications
2. **Real Blockchain**: Connect to actual blockchain networks
3. **Advanced Evolution**: Implement real genetic algorithms
4. **Machine Learning**: Add actual ML-based contract optimization
5. **Collaboration Networks**: Build inter-contract communication

## 🛠️ Troubleshooting

### Server Not Running
```bash
# Start the server
go run ./cmd/server

# Check if port 8080 is available
lsof -i :8080
```

### Dependencies Issues
```bash
# Download dependencies
go mod download

# Clean module cache
go clean -modcache
go mod download
```

### API Testing Issues
```bash
# Test basic connectivity
curl http://localhost:8080/health

# Check server logs for errors
# (Server logs appear in terminal where you started it)
```

---

**🧬 Living Smart Contracts API - Revolutionary blockchain platform with evolutionary intelligence!**
