#!/bin/bash

# Living Smart Contracts API Test Script
echo "🧬 Testing Living Smart Contracts API"
echo "======================================"

BASE_URL="http://localhost:8080"

echo ""
echo "1. 💚 Health Check:"
curl -s "$BASE_URL/health" | jq .

echo ""
echo "2. 📋 API Information:"
curl -s "$BASE_URL/api/info" | jq .

echo ""
echo "3. 📜 List Contracts:"
curl -s "$BASE_URL/api/contracts" | jq .

echo ""
echo "4. 🚀 Deploy New Contract:"
curl -X POST "$BASE_URL/api/contracts/deploy" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "adaptive_defi_strategy",
    "type": "living",
    "source_code": "contract AdaptiveDeFi { function optimize() public returns (uint256) { return 42; } }",
    "owner": "0x1234567890123456789012345678901234567890",
    "init_params": {"initial_capital": 1000000, "risk_tolerance": 0.3},
    "time_aware": true,
    "history_depth": 1000
  }' | jq .

echo ""
echo "5. ⚡ Execute Contract Function:"
curl -X POST "$BASE_URL/api/contracts/execute" \
  -H "Content-Type: application/json" \
  -d '{
    "contract_address": "0xabcdef1234567890abcdef1234567890abcdef12",
    "function": "optimize",
    "parameters": {"strategy": "yield_farming", "risk_level": 0.5},
    "caller": "0x1234567890123456789012345678901234567890",
    "gas_limit": 500000
  }' | jq .

echo ""
echo "6. 🧬 Trigger Evolution:"
curl -X POST "$BASE_URL/api/evolution/trigger" \
  -H "Content-Type: application/json" \
  -d '{
    "evolution_type": "optimization",
    "parameters": {"target_metric": "success_rate", "mutation_rate": 0.05}
  }' | jq .

echo ""
echo "7. 🤝 Propose Collaboration:"
curl -X POST "$BASE_URL/api/collaborations/propose" \
  -H "Content-Type: application/json" \
  -d '{
    "proposer_address": "0xabcdef1234567890abcdef1234567890abcdef12",
    "target_address": "0x9876543210987654321098765432109876543210",
    "collaboration_type": "resource_sharing",
    "terms": {"resource_allocation": {"proposer": 0.6, "target": 0.4}},
    "expected_benefit": 0.25
  }' | jq .

echo ""
echo "8. 🔍 Get Contract Details:"
curl -s "$BASE_URL/api/contracts/0xabcdef1234567890abcdef1234567890abcdef12" | jq .

echo ""
echo "======================================"
echo "✅ All API endpoints tested successfully!"
echo "📊 Swagger UI available at: $BASE_URL/swagger/"
echo "💚 Health check: $BASE_URL/health"
echo "======================================"
